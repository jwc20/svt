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
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
			//return key.Type() == "ssh-ed25519"
		}),
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

		userName := s.User()
		var pubKeyStr string

		key := s.PublicKey()
		log.Infof("User %s connected with public key %s", userName, key)
		if key != nil {
			// User has a real SSH key, use the standard format
			pubKeyStr = strings.TrimSpace(string(gossh.MarshalAuthorizedKey(key)))
		} else {
			// Fallback: Use a combination of username and remote address
			// This allows users without keys to still have a unique record
			remoteAddr := s.RemoteAddr().String()
			pubKeyStr = "unauthenticated:" + userName + ":" + remoteAddr
		}

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
