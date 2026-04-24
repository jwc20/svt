package engine

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Serialize encodes the game state and turn history into a compact string.
// Format: [turn];[infra];[turnHistory];[currentState] - [eventHistory]
//
// Example: 3;s4d3;a1/a2w/a1;1200/65/4/2/420/650 - 0/3/15
func Serialize(gs *GameState) string {
	// Turn number
	turn := strconv.Itoa(gs.TurnNumber)

	// Infrastructure: s{1-4}d{1-3}
	// ServerOption iota: Fargate=0, EC2=1, Lambda=2, ThinkPad=3 → +1 for 1-indexed
	// DBOption iota: Aurora=0, RDS=1, SQLite=2 → +1 for 1-indexed
	infra := fmt.Sprintf("s%dd%d", int(gs.Server)+1, int(gs.Database)+1)

	// Turn history: a1/a2w/a2l/...
	turns := make([]string, len(gs.TurnHistory))
	for i, t := range gs.TurnHistory {
		s := fmt.Sprintf("a%d", t.Action)
		switch t.Action {
		case 2:
			if t.DeathRoll == DeathRollWin {
				s += "w"
			} else {
				s += "l"
			}
		case 3:
			// no suffix needed for marketing push
		}
		turns[i] = s
	}
	turnHistory := strings.Join(turns, "/")

	// Current state: cash/hype/techDebt/bugCount/miles/userCount
	state := fmt.Sprintf("%d/%d/%d/%d/%d/%d",
		gs.Cash, gs.Hype, gs.TechDebt, gs.BugCount, gs.ProductReadiness, gs.UserCount)

	// Event history: eventID per turn
	events := make([]string, len(gs.TurnHistory))
	for i, t := range gs.TurnHistory {
		events[i] = strconv.Itoa(t.EventID)
	}
	eventHistory := strings.Join(events, "/")

	return fmt.Sprintf("%s;%s;%s;%s - %s", turn, infra, turnHistory, state, eventHistory)
}

// Deserialize parses a compact state string back into a GameState.
func Deserialize(s string) (*GameState, error) {
	// Split on " - " to separate state from event history
	halves := strings.SplitN(s, " - ", 2)
	if len(halves) != 2 {
		return nil, errors.New("invalid format: missing ' - ' separator")
	}

	left := halves[0]
	eventPart := halves[1]

	// Left side: turn;infra;turnHistory;currentState
	sections := strings.SplitN(left, ";", 4)
	if len(sections) != 4 {
		return nil, errors.New("invalid format: expected 4 sections before ' - '")
	}

	// Parse turn number
	turnNumber, err := strconv.Atoi(sections[0])
	if err != nil {
		return nil, fmt.Errorf("invalid turn number: %w", err)
	}

	// Parse infra: s{N}d{N}
	infraStr := sections[1]
	if len(infraStr) < 4 || infraStr[0] != 's' || !strings.Contains(infraStr, "d") {
		return nil, fmt.Errorf("invalid infra format: %s", infraStr)
	}
	dIdx := strings.Index(infraStr, "d")
	serverNum, err := strconv.Atoi(infraStr[1:dIdx])
	if err != nil || serverNum < 1 || serverNum > 4 {
		return nil, fmt.Errorf("invalid server choice: %s", infraStr)
	}
	dbNum, err := strconv.Atoi(infraStr[dIdx+1:])
	if err != nil || dbNum < 1 || dbNum > 3 {
		return nil, fmt.Errorf("invalid db choice: %s", infraStr)
	}

	// Parse current state: cash/hype/techDebt/bugCount/miles/userCount
	stateParts := strings.Split(sections[3], "/")
	if len(stateParts) != 6 {
		return nil, errors.New("invalid state: expected 6 values")
	}
	stateVals := make([]int, 6)
	for i, p := range stateParts {
		v, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid state value at index %d: %w", i, err)
		}
		stateVals[i] = v
	}

	// Parse turn history and event history
	var turnHistory []TurnEntry

	if sections[2] != "" {
		turnParts := strings.Split(sections[2], "/")
		eventParts := strings.Split(eventPart, "/")
		if len(turnParts) != len(eventParts) {
			return nil, fmt.Errorf("turn history (%d) and event history (%d) length mismatch",
				len(turnParts), len(eventParts))
		}

		turnHistory = make([]TurnEntry, len(turnParts))
		for i, tp := range turnParts {
			if len(tp) < 2 || tp[0] != 'a' {
				return nil, fmt.Errorf("invalid turn entry: %s", tp)
			}

			entry := TurnEntry{}
			actionChar := tp[1]
			switch actionChar {
			case '1':
				entry.Action = 1
				entry.DeathRoll = DeathRollNone
			case '2':
				entry.Action = 2
				if len(tp) < 3 {
					return nil, fmt.Errorf("action 2 missing death roll result: %s", tp)
				}
				switch tp[2] {
				case 'w':
					entry.DeathRoll = DeathRollWin
				case 'l':
					entry.DeathRoll = DeathRollLoss
				default:
					return nil, fmt.Errorf("invalid death roll result: %s", tp)
				}
			case '3':
				entry.Action = 3
				entry.DeathRoll = DeathRollNone
			default:
				return nil, fmt.Errorf("invalid action: %s", tp)
			}

			eid, err := strconv.Atoi(eventParts[i])
			if err != nil {
				return nil, fmt.Errorf("invalid event ID at index %d: %w", i, err)
			}
			entry.EventID = eid

			turnHistory[i] = entry
		}
	}

	gs := &GameState{
		TurnNumber:       turnNumber,
		Server:           ServerOption(serverNum - 1),
		Database:         DBOption(dbNum - 1),
		Cash:             stateVals[0],
		Hype:             stateVals[1],
		TechDebt:         stateVals[2],
		BugCount:         stateVals[3],
		ProductReadiness: stateVals[4],
		UserCount:        stateVals[5],
		TurnHistory:      turnHistory,
	}

	return gs, nil
}
