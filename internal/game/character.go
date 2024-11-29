package game

import (
	"fmt"
)

// Character represents a playable character in the battle card game
type Character struct {
	ID            string           `json:"ID"`
	Name          string           `json:"Name"`
	Abilities     []Ability        `json:"Abilities"`
	StatusEffects []StatusEffectData `json:"StatusEffects"`
	Health        int              `json:"Health"`
	Attack        int              `json:"Attack"`
	Defense       int              `json:"Defense"`
	Speed         int              `json:"Speed"`
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

	// Calculate total damage based on ability damage and character's Attack stat
	damage := ability.Damage + c.Attack

	// Apply the damage to target
	target.TakeDamage(damage)

	// Apply status effect if present
	if ability.StatusEffect.Type != "" {
		target.StatusEffects = append(target.StatusEffects, ability.StatusEffect)
	}

	// Call ability.Use() to set cooldown
	ability.Use()

	// Return result with information about what happened
	return AbilityResult{
		Success:      true,
		Damage:       damage,
		StatusEffect: &ability.StatusEffect,
		Message:      fmt.Sprintf("Ability used successfully for %d damage", damage),
	}
}

// Process status effect - handle all active status effects.
// loop over each effect in the StatusEffectData slice
// switch fof each type of StatusEffect
//
func (c *Character) ProcessStatusEffect() {
	// Create a new slice to store active effects
	activeEffects := make([]StatusEffectData, 0)

	// Process each status effect
	for _, effect := range c.StatusEffects {
		if effect.Duration <= 0 {
			continue
		}

		// Process effect based on type
		switch effect.Type {
		case StatusAccelerate:
			// Increase speed based on potency
			speedIncrease := c.GetEffectScalingValue("speed", 100, effect.Potency, 1)
			c.Speed = speedIncrease

		case StatusBurning:
			// Increasing damage over time
			burnDamage := c.GetEffectScalingValue("health", 100, effect.Potency*effect.Duration, 10)
			c.TakeDamage(burnDamage)

		case StatusPoisoned:
			// Decreasing damage over time
			poisonDamage := c.GetEffectScalingValue("health", 100, effect.Potency, effect.Duration)
			c.TakeDamage(poisonDamage)

		case StatusEnraged:
			// Increase attack based on potency
			attackIncrease := c.GetEffectScalingValue("attack", 100, effect.Potency, 1)
			c.Attack = attackIncrease

		case StatusRegenerating:
			// Heal based on potency
			healAmount := c.GetEffectScalingValue("health", 100, effect.Potency, 10)
			c.Health += healAmount
		}

		// Decrease duration and keep active effects
		effect.Duration--
		if effect.Duration > 0 {
			activeEffects = append(activeEffects, effect)
		}
	}

	// Update status effects list with only active effects
	c.StatusEffects = activeEffects
}

// Get Effect Value
func (c *Character) GetEffectScalingValue(statName string, scalar int, potency int, divisor int) int {
	// Get base stat value based on statName
	var baseStat int
	switch statName {
	case "health":
		baseStat = c.Health
	case "attack":
		baseStat = c.Attack
	case "defense":
		baseStat = c.Defense
	case "speed":
		baseStat = c.Speed
	default:
		fmt.Printf("Warning: unknown stat %s\n", statName)
		return 0
	}
	// Calculate stat modification
	modifiedStat := baseStat + (baseStat * potency / scalar)
	// Round down and apply divisor
	newStat := int(modifiedStat / divisor)

	fmt.Printf("newStat: %d\n", newStat)
	return newStat
}
