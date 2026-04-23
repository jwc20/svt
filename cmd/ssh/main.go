package main

import (
	"bufio"
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

// buildMode is set at build time via ldflags: "prod" or "local"
var buildMode = "local"

func loadEnv(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
	}
}

func getHostPort() (string, string) {
	loadEnv(".env")

	var host, port string
	if buildMode == "prod" {
		host = os.Getenv("PROD_HOST")
		port = os.Getenv("PROD_PORT")
	} else {
		host = os.Getenv("LOCAL_HOST")
		port = os.Getenv("LOCAL_PORT")
	}

	// Defaults
	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = "22"
	}
	return host, port
}

func main() {
	host, port := getHostPort()

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
			wish.Fatalf(s, "no public key found, skipping")
			return nil
		}

		playerID, err := db.CreatePlayer(pubKeyStr, userName)
		if err != nil {
			wish.Fatalf(s, "could not create player: %v", err)
			return nil
		}

		bonusHype, err := db.GetBonusHype(playerID)
		if err != nil || bonusHype < 0 {
			bonusHype = hackernews.FetchBonusHype(userName)
			if err := db.SetBonusHype(playerID, bonusHype); err != nil {
				log.Error("could not persist bonus hype",
					"player_id", playerID,
					"user", userName,
					"error", err,
				)
			}
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
