package game

// Character represents a playable character in the battle card game
type Character struct {
	Name          string
	Abilities     []Ability
	StatusEffects []StatusEffectData
	ID            int
	Health        int
	Attack        int
	Defense       int
	Speed         int
}

// IsValid will check if the character has valid stats
func (c Character) IsValid() bool {
	return c.Name != "" &&
		c.Health > 0 &&
		c.Attack >= 0 &&
		c.Defense >= 0 &&
		c.Speed > 0
}

func (c *Character) TakeDamage(damage int) {
	actualDamage := damage - c.Defense
	if actualDamage < 0 {
		actualDamage = 0
	}
	c.Health -= actualDamage
	if c.Health < 0 {
		c.Health = 0
	}
}

// Will allow character to use ability on a target
func (c *Character) UseAbility(abilityIndex int, target *Character) AbilityResult {
	// Make sure the index is within the bounds of the []Abilities
	if abilityIndex >= len(c.Abilities) {
		return AbilityResult{
			Success: false,
			Message: "Invalid ability index.",
		}
	}
	// Check if ability can be used
	ability := &c.Abilities[abilityIndex]
	if !ability.CanUse() {
		return AbilityResult{
			Success: false,
			Message: "Ability on cooldown",
		}
	}
	// Apply the damag
	// Apply status effect
	damage := ability.Damage
	target.TakeDamage(damage)

	if ability.StatusEffect.Type != "" {
		target.StatusEffects = append(target.StatusEffects, ability.StatusEffect)
	}

	// Call ability.Use() to set cooldown
	ability.Use()
	// Return result: Information about what happened.
	return AbilityResult{
		Success:      true,
		Damage:       damage,
		StatusEffect: &ability.StatusEffect,
		Message:      "Ability used successfully",
	}
}

// Process status effect - handle all active status effects.
