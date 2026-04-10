package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/log/v2"
	"charm.land/wish/v2"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/charmbracelet/ssh"

	trail "github.com/jwc20/svt"
)

const (
	host = "localhost"
	port = "23234"
)

type SimpleStore struct{}

func (s *SimpleStore) SaveState(state trail.GameState) error { return nil }
func (s *SimpleStore) LoadState() (trail.GameState, error)   { return trail.GameState{}, nil }

func main() {
	s, err := wish.NewServer(
		ssh.AllocatePty(),
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			SvtBubbleteaMiddleware(),
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

func SvtBubbleteaMiddleware() wish.Middleware {
	store := &SimpleStore{}

	teaHandler := func(s ssh.Session) *tea.Program {
		pty, _, active := s.Pty()
		if !active {
			wish.Fatalln(s, "no active terminal, skipping")
			return nil
		}

		playerId := s.User()

		m := trail.NewRootModel(store, playerId)
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
