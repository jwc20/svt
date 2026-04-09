package svt

type StubGameStore struct {
	State GameState
	Saved bool
}

func (s *StubGameStore) SaveState(state GameState) error {
	s.Saved = true
	s.State = state
	return nil
}

func (s *StubGameStore) LoadState() (GameState, error) {
	return s.State, nil
}
