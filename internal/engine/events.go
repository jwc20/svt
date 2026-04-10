package engine

import (
	"github.com/jwc20/svt/internal/rand"
)

func GenerateEvent(gs *GameState) string {
	r := rand.GetRandomInt(100)

	switch {
	case r <= 6:
		gs.Trip.Mileage -= 15 + rand.GetRandomInt(10)
		gs.Inventory.Miscellaneous -= 8 + rand.GetRandomInt(5)
		return "WAGON BREAKS DOWN — LOSS OF TIME AND SUPPLIES."
	case r <= 11:
		gs.Trip.Mileage -= 25
		gs.Inventory.Oxen -= 20
		return "OX INJURED — LOSS OF TIME."
	case r <= 15:
		gs.Flags.Injured = true
		gs.Inventory.Miscellaneous -= 5 + rand.GetRandomInt(4)
		return "BAD LUCK — YOUR DAUGHTER BROKE HER ARM."
	case r <= 20:
		gs.Inventory.Ammo -= 10 + rand.GetRandomInt(5)
		if gs.Inventory.Ammo < 0 {
			gs.Inventory.Food -= 30 + rand.GetRandomInt(20)
			return "WILD ANIMALS ATTACK! YOU RAN OUT OF BULLETS — THEY GOT SOME FOOD."
		}
		return "WILD ANIMALS ATTACK!"
	case r <= 25:
		if gs.Inventory.Clothing < 20 {
			gs.Flags.Ill = true
			return "COLD WEATHER — BRRRR! YOU DON'T HAVE ENOUGH CLOTHING."
		}
		return "COLD WEATHER — BRRRR! BUT YOU'RE DRESSED WARM."
	case r <= 30:
		gs.Trip.Mileage -= 10 + rand.GetRandomInt(5)
		gs.Inventory.Food -= 10
		gs.Inventory.Ammo -= 5 + rand.GetRandomInt(5)
		gs.Inventory.Miscellaneous -= 5 + rand.GetRandomInt(5)
		return "HEAVY RAINS — TIME LOST AND SUPPLIES DAMAGED."
	case r <= 33:
		gs.Inventory.Food -= 10 + rand.GetRandomInt(10)
		gs.Player.Cash -= 10 + rand.GetRandomInt(15)
		if gs.Player.Cash < 0 {
			gs.Player.Cash = 0
		}
		return "BANDITS ATTACK!"
	case r <= 36:
		gs.Inventory.Food -= 40 + rand.GetRandomInt(30)
		gs.Inventory.Ammo -= 20 + rand.GetRandomInt(20)
		gs.Inventory.Miscellaneous -= 10 + rand.GetRandomInt(10)
		return "FIRE IN YOUR WAGON — LOSS OF SUPPLIES."
	case r <= 40:
		gs.Inventory.Food += 14 + rand.GetRandomInt(5)
		return "HELPFUL INDIANS SHOW YOU WHERE TO FIND FOOD."
	default:
		return ""
	}
}

func HandleAilment(gs *GameState) (bool, string) {
	if gs.Inventory.Miscellaneous < 5 {
		if gs.Flags.Ill {
			return false, "YOU DIED OF PNEUMONIA."
		}
		return false, "YOU DIED OF YOUR INJURIES."
	}
	gs.Inventory.Miscellaneous -= 5 + rand.GetRandomInt(5)
	gs.Flags.Ill = false
	gs.Flags.Injured = false
	return true, "YOU USED MEDICINE AND RESTED."
}
