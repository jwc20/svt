package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"charm.land/wish/v2"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/charmbracelet/ssh"
	gossh "golang.org/x/crypto/ssh"

	"github.com/jwc20/svt/internal/hackernews"
	"github.com/jwc20/svt/internal/store"
	"github.com/jwc20/svt/internal/ui"
)

const (
	host = "0.0.0.0"
	port = "22"
)

func main() {
	db, err := store.NewSQLiteStore("svt.db")
	if err != nil {
		log.Fatal("Could not open database", "error", err)
	}
	defer db.Close()

	s, err := wish.NewServer(
		ssh.AllocatePty(),
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		//wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		//	return key.Type() == "ssh-ed25519"
		//}),
		wish.WithMiddleware(
			SvtBubbleteaMiddleware(db),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func SvtBubbleteaMiddleware(db *store.SQLiteStore) wish.Middleware {
	teaHandler := func(s ssh.Session) *tea.Program {
		pty, _, active := s.Pty()
		if !active {
			wish.Fatalln(s, "no active terminal, skipping")
			return nil
		}

		key := s.PublicKey()
		if key == nil {
			wish.Fatalln(s, "public key required")
			return nil

			// https://github.com/charmbracelet/ssh/blob/ebfa259c73091350caed965eb59c2bb8cd90e7e1/_examples/ssh-publickey/public_key.go#L14
			// parsed, _, _, _, _ := ssh.ParseAuthorizedKey(
			// 	[]byte("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAILxWe2rXKoiO6W14LYPVfJKzRfJ1f3Jhzxrgjc/D4tU7"),
			// )
			// key = parsed
		}

		// Marshal the public key to authorized_keys format (e.g. "ssh-ed25519 AAAA...")
		pubKeyStr := strings.TrimSpace(string(gossh.MarshalAuthorizedKey(key)))

		userName := s.User()

		playerID, err := db.CreatePlayer(pubKeyStr, userName)
		if err != nil {
			wish.Fatalf(s, "could not create player: %v", err)
			return nil
		}

		bonusHype, err := db.GetBonusHype(playerID)
		if err != nil || bonusHype < 0 {
			bonusHype = hackernews.FetchBonusHype(userName)
			_ = db.SetBonusHype(playerID, bonusHype)
		}

		m := ui.NewRootModel(db, playerID, userName, bonusHype)
		opts := bubbletea.MakeOptions(s)

		p := tea.NewProgram(m, opts...)
		go func() {
			time.Sleep(100 * time.Millisecond)
			p.Send(tea.WindowSizeMsg{Width: pty.Window.Width, Height: pty.Window.Height})
		}()
		return p
	}
	return bubbletea.MiddlewareWithProgramHandler(teaHandler)
}
