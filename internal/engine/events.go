package engine

import (
	"github.com/jwc20/svt/internal/rand"
)

type EventEffect struct {
	CashDelta     int
	HypeDelta     int
	MilesDelta    int
	TechDebtDelta int
	BugDelta      int
}

type Event struct {
	Name   string
	Effect EventEffect
}

// Event IDs — stable integers for serialization. Never reorder, only append.
const (
	EventNone                = 0
	EventDNSOutage           = 1
	EventLeadEngineerQuits   = 2
	EventLaptopStolen        = 3
	EventDDoS                = 4
	EventAWSBillSurprise     = 5
	EventDatabaseCorruption  = 6
	EventCompetitorLaunches  = 7
	EventSecurityBreach      = 8
	EventNpmBreaks           = 9
	EventInternPushes        = 10
	EventRentRaised          = 11
	EventRansomware          = 12
	EventTechBlogReview      = 13
	EventHackerNewsViral     = 14
	EventVCColdEmail         = 15
	EventRichUncle           = 16
	EventOpenSourceFix       = 17
	EventPitchCompetition    = 18
	EventInfluencerTweet     = 19
)

// GenerateEvent rolls a random event and applies it to the game state.
// Returns the event description and the event ID (0 if nothing happens).
func GenerateEvent(gs *GameState) (string, int) {
	r := rand.GetRandomInt(101) - 1 // 0..100

	var evt Event
	var eventID int

	switch {
	case r <= 5:
		eventID = EventDNSOutage
		evt = Event{
			Name: "DNS outage -- site unreachable for hours.",
			Effect: EventEffect{
				MilesDelta: -(15 + rand.GetRandomInt(11) - 1),
				HypeDelta:  -(rand.GetRandomInt(6) + 4), // rand(5..10)
			},
		}
	case r <= 10:
		eventID = EventLeadEngineerQuits
		evt = Event{
			Name: "Lead engineer quits.",
			Effect: EventEffect{
				TechDebtDelta: 10 + rand.GetRandomInt(6) - 1, // 10 + rand(0..5)
				BugDelta:      rand.GetRandomInt(3),           // rand(1..3)
			},
		}
	case r <= 14:
		eventID = EventLaptopStolen
		evt = Event{
			Name: "Cofounder's laptop stolen at coffee shop.",
			Effect: EventEffect{
				CashDelta:     -(100 + rand.GetRandomInt(76) - 1), // cash -= 100 + rand(0..75)
				TechDebtDelta: rand.GetRandomInt(4) + 2,           // rand(3..6)
			},
		}
	case r <= 19:
		eventID = EventDDoS
		evt = Event{
			Name: "DDoS attack!",
			Effect: EventEffect{
				MilesDelta: -25,
				HypeDelta:  -(rand.GetRandomInt(11) + 9), // rand(10..20)
			},
		}
	case r <= 24:
		eventID = EventAWSBillSurprise
		evt = Event{
			Name: "AWS bill surprise -- forgot to turn off test instances.",
			Effect: EventEffect{
				CashDelta: -(150 + rand.GetRandomInt(101) - 1), // cash -= 150 + rand(0..100)
			},
		}
	case r <= 28:
		eventID = EventDatabaseCorruption
		evt = Event{
			Name: "Database corruption -- loss of user data.",
			Effect: EventEffect{
				HypeDelta: -(15 + rand.GetRandomInt(16) - 1), // hype -= 15 + rand(0..15)
				BugDelta:  rand.GetRandomInt(3) + 1,           // rand(2..4)
			},
		}
	case r <= 32:
		eventID = EventCompetitorLaunches
		evt = Event{
			Name: "Competitor launches same feature. Users jump ship.",
			Effect: EventEffect{
				HypeDelta:  -(10 + rand.GetRandomInt(11) - 1), // hype -= 10 + rand(0..10)
				MilesDelta: -(rand.GetRandomInt(11) + 4),      // miles -= rand(5..15)
			},
		}
	case r <= 35:
		eventID = EventSecurityBreach
		evt = Event{
			Name: "Security breach -- passwords leaked.",
			Effect: EventEffect{
				HypeDelta:     -(20 + rand.GetRandomInt(11) - 1), // hype -= 20 + rand(0..10)
				CashDelta:     -(50 + rand.GetRandomInt(51) - 1), // cash -= 50 + rand(0..50)
				TechDebtDelta: rand.GetRandomInt(3) + 2,          // rand(3..5)
			},
		}
	case r <= 38:
		eventID = EventNpmBreaks
		evt = Event{
			Name: "npm dependency breaks -- half the app crashes.",
			Effect: EventEffect{
				BugDelta:      5 + rand.GetRandomInt(4) - 1, // 5 + rand(0..3)
				TechDebtDelta: rand.GetRandomInt(4) + 1,     // rand(2..5)
			},
		}
	case r <= 41:
		eventID = EventInternPushes
		evt = Event{
			Name: "Intern pushes to prod on Friday night.",
			Effect: EventEffect{
				BugDelta:  rand.GetRandomInt(4) + 2, // rand(3..6)
				HypeDelta: -(rand.GetRandomInt(6) + 4), // rand(5..10)
			},
		}
	case r <= 44:
		eventID = EventRentRaised
		evt = Event{
			Name: "Office landlord raises rent.",
			Effect: EventEffect{
				CashDelta: -(75 + rand.GetRandomInt(51) - 1), // cash -= 75 + rand(0..50)
			},
		}
	case r <= 49:
		eventID = EventRansomware
		if gs.TechDebt < 30 {
			evt = Event{
				Name: "Ransomware on a dev machine.",
				Effect: EventEffect{
					CashDelta: -(250 + rand.GetRandomInt(151) - 1), // cash -= 250 + rand(0..150)
				},
			}
		} else {
			evt = Event{
				Name: "Ransomware on a dev machine.",
				Effect: EventEffect{
					TechDebtDelta: rand.GetRandomInt(3) + 1, // rand(2..4)
				},
			}
		}
	case r <= 54:
		eventID = EventTechBlogReview
		evt = Event{
			Name: "Tech blog writes a positive review!",
			Effect: EventEffect{
				HypeDelta: 10 + rand.GetRandomInt(11) - 1, // hype += 10 + rand(0..10)
			},
		}
	case r <= 58:
		eventID = EventHackerNewsViral
		evt = Event{
			Name: "Post goes viral on Hacker News.",
			Effect: EventEffect{
				HypeDelta:  15 + rand.GetRandomInt(16) - 1, // hype += 15 + rand(0..15)
				MilesDelta: rand.GetRandomInt(11) + 4,      // miles += rand(5..15)
			},
		}
	case r <= 62:
		eventID = EventVCColdEmail
		evt = Event{
			Name: "VC cold-emails you after seeing your product.",
			Effect: EventEffect{
				CashDelta: 250 + rand.GetRandomInt(251) - 1, // cash += 250 + rand(0..250)
			},
		}
	case r <= 65:
		eventID = EventRichUncle
		evt = Event{
			Name: "Rich uncle sends a check.",
			Effect: EventEffect{
				CashDelta: 150 + rand.GetRandomInt(101) - 1, // cash += 150 + rand(0..100)
			},
		}
	case r <= 68:
		eventID = EventOpenSourceFix
		evt = Event{
			Name: "Open source contributor fixes 3 bugs for free.",
			Effect: EventEffect{
				BugDelta:      -3,
				TechDebtDelta: -(rand.GetRandomInt(4) + 1), // techDebt -= rand(2..5)
			},
		}
	case r <= 71:
		eventID = EventPitchCompetition
		evt = Event{
			Name: "Win a startup pitch competition.",
			Effect: EventEffect{
				CashDelta: 200 + rand.GetRandomInt(101) - 1, // cash += 200 + rand(0..100)
				HypeDelta: rand.GetRandomInt(6) + 4,         // hype += rand(5..10)
			},
		}
	case r <= 74:
		eventID = EventInfluencerTweet
		evt = Event{
			Name: "Influencer tweets about your product.",
			Effect: EventEffect{
				HypeDelta: 10 + rand.GetRandomInt(11) - 1, // hype += 10 + rand(0..10)
			},
		}
	default:
		// 75-100: Nothing happens
		return "", EventNone
	}

	// Apply effects
	gs.Cash += evt.Effect.CashDelta
	gs.Hype += evt.Effect.HypeDelta
	if gs.Hype < 0 {
		gs.Hype = 0
	}
	gs.Miles += evt.Effect.MilesDelta
	if gs.Miles < 0 {
		gs.Miles = 0
	}
	gs.TechDebt += evt.Effect.TechDebtDelta
	if gs.TechDebt < 0 {
		gs.TechDebt = 0
	}
	gs.BugCount += evt.Effect.BugDelta
	if gs.BugCount < 0 {
		gs.BugCount = 0
	}

	return evt.Name, eventID
}

// CheckIncident checks if an incident occurs and applies consequences.
// Returns (survived, message).
func CheckIncident(gs *GameState) (bool, string) {
	th := TechHealth(gs)
	threshold := 10 + rand.GetRandomInt(6) - 1 // 10 + rand(0..5)

	if th > threshold {
		return true, ""
	}

	// Incident! Apply penalties
	gs.Hype -= rand.GetRandomInt(21) + 9 // hype -= rand(10..30)
	if gs.Hype < 0 {
		gs.Hype = 0
	}
	gs.BugCount += rand.GetRandomInt(3) // bugCount += rand(1..3)

	return false, "SYSTEM INCIDENT! Infrastructure is failing under the weight of tech debt and bugs."
}
