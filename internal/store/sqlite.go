package store

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/jwc20/svt/internal/engine"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	s := &SQLiteStore{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return s, nil
}

func (s *SQLiteStore) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS players (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			public_key TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS games (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			state TEXT NOT NULL,
			player_id INTEGER NOT NULL REFERENCES players(id),
			score INTEGER,
			ended BOOLEAN NOT NULL DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_games_player_id ON games(player_id);
	`)
	if err != nil {
		return err
	}

	// Add ended column to existing databases that lack it.
	_, err = s.db.Exec(`ALTER TABLE games ADD COLUMN ended BOOLEAN NOT NULL DEFAULT FALSE`)
	if err != nil {
		// Ignore "duplicate column" error — column already exists.
		if !strings.Contains(err.Error(), "duplicate column") {
			return err
		}
	}

	// Backfill: mark games that have a score as ended.
	_, err = s.db.Exec(`UPDATE games SET ended = TRUE WHERE score IS NOT NULL AND ended = FALSE`)
	return err
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) CreatePlayer(publicKey string) (int64, error) {
	result, err := s.db.Exec(
		`INSERT OR IGNORE INTO players (public_key) VALUES (?)`,
		publicKey,
	)
	if err != nil {
		return 0, err
	}

	// If the row was inserted, return the new ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	if id > 0 {
		return id, nil
	}

	// Row already existed, look up the ID
	return s.GetPlayerByKey(publicKey)
}

func (s *SQLiteStore) GetPlayerByKey(publicKey string) (int64, error) {
	var id int64
	err := s.db.QueryRow(
		`SELECT id FROM players WHERE public_key = ?`,
		publicKey,
	).Scan(&id)
	return id, err
}

func (s *SQLiteStore) NewGame(playerID int64, state *engine.GameState) (int64, error) {
	serialized := engine.Serialize(state)
	result, err := s.db.Exec(`
		INSERT INTO games (state, player_id)
		VALUES (?, ?)
	`, serialized, playerID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *SQLiteStore) SaveGame(playerID int64, gameID int64, state *engine.GameState) error {
	serialized := engine.Serialize(state)
	_, err := s.db.Exec(`
		UPDATE games SET state = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND player_id = ?
	`, serialized, gameID, playerID)
	return err
}

func (s *SQLiteStore) LoadActiveGame(playerID int64) (int64, *engine.GameState, error) {
	var gameID int64
	var stateStr string
	err := s.db.QueryRow(`
		SELECT id, state FROM games
		WHERE player_id = ? AND ended = FALSE
		ORDER BY updated_at DESC
		LIMIT 1
	`, playerID).Scan(&gameID, &stateStr)

	if err == sql.ErrNoRows {
		return 0, nil, nil
	}
	if err != nil {
		return 0, nil, err
	}

	gs, err := engine.Deserialize(stateStr)
	if err != nil {
		return 0, nil, fmt.Errorf("deserialize: %w", err)
	}

	return gameID, gs, nil
}

func (s *SQLiteStore) FinishGame(gameID int64, score *int) error {
	_, err := s.db.Exec(`
		UPDATE games SET score = ?, ended = TRUE, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, score, gameID)
	return err
}

func (s *SQLiteStore) Leaderboard(limit int) ([]engine.LeaderboardEntry, error) {
	rows, err := s.db.Query(`
		SELECT p.public_key, g.score, g.updated_at
		FROM games g
		JOIN players p ON p.id = g.player_id
		WHERE g.ended = TRUE AND g.score IS NOT NULL
		ORDER BY g.score DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []engine.LeaderboardEntry
	rank := 1
	for rows.Next() {
		var e engine.LeaderboardEntry
		if err := rows.Scan(&e.PublicKey, &e.Score, &e.EndedAt); err != nil {
			return nil, err
		}
		e.Rank = rank
		rank++
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
