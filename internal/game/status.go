package game

// Status effects are strings and there are many types
type StatusEffect string

const (
	StatusAccelerate   StatusEffect = "ACCELERATE" // Each round increase speed by percentage based on potency.
	StatusBurning      StatusEffect = "BURNING"    // Damage over time, damage increases by some formula that uses the duration each round.
	StatusPoisoned     StatusEffect = "POISON"     // Damage over time, damage decreases by some formula that uses the duration each time.
	StatusEnraged      StatusEffect = "ENRAGED"    // Each round increase Attack power by percentage based on potency.
	StatusRegenerating StatusEffect = "REGENERATING"
)

// EVery effect will hold various attributes
type StatusEffectData struct {
	Type     StatusEffect `json:"Type"`
	Duration int         `json:"Duration"` // Number of remaining turns
	Potency  int         `json:"Potency"`  // The strength of the effect
}
