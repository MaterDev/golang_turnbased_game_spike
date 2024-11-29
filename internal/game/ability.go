package game

// Ability represents a unique move of a character. Characters can have multiple abilities.
type Ability struct {
	Name         string          `json:"Name"`
	StatusEffect StatusEffectData `json:"StatusEffect"`
	Damage       int             `json:"Damage"`
	CooldownMax  int             `json:"CooldownMax"`
	Cooldown     int             `json:"Cooldown"`
}

// AbilityResult contains the result of using an ability
type AbilityResult struct {
	Success      bool             `json:"Success"`
	Damage       int              `json:"Damage"`
	StatusEffect *StatusEffectData `json:"StatusEffect,omitempty"`
	Message      string           `json:"Message"`
}

func (a *Ability) CanUse() bool {
	// Will only return true if cooldown is 0 to prevent overuse of ability
	return a.Cooldown == 0
}

func (a *Ability) Use() {
	if a.CanUse() {
		// When ability is used, will set the cooldown.
		a.Cooldown = a.CooldownMax
	}
}

func (a *Ability) ReduceCooldown() {
	if a.Cooldown > 0 {
		a.Cooldown--
	}
}
