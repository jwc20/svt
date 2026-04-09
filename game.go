package svt

const TotalRequiredMileage = 2040
const InitialCash = 700

type Player struct {
	Cash          int
	ShootingLevel int
}

type Inventory struct {
	Oxen          int
	Food          int
	Ammo          int
	Clothing      int
	Miscellaneous int
}
