package game

// Character represents a playable character in the battle card game
type Character struct {
	ID          int
	Name        string
	Health      int
	Attack      int
	Defence     int
	Speed       int
	SpecialMove string
}

// IsValid will check if the character has valid stats
func (c Character) IsValid() bool {
	return c.Name != "" &&
		c.Health > 0 &&
		c.Attack >= 0 &&
		c.Defence >= 0 &&
		c.Speed > 0
}
