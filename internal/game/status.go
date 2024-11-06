package game

// Status effects are strings and there are many types
type StatusEffect string

const (
	StatusStunned      StatusEffect = "STUNNED"
	StatusBurning      StatusEffect = "BURNING"
	StatusPoisoned     StatusEffect = "POISON"
	StatusShielded     StatusEffect = "SHIELDED"
	StatusRegenerating StatusEffect = "REGENERATING"
)

// EVery effect will hold various attributes
type StatusEffectData struct {
	Type     StatusEffect
	Duration int // Number of remaining turns
	Potency  int // The strength of the effect
}
