package game

// Ability represents a unique move of a character. Characters can have multiple abilities.
type Ability struct {
	Name         string
	StatusEffect StatusEffectData
	Damage       int
	CooldownMax  int
	Cooldown     int
}

// AbilityResult contains the result of using an ability
type AbilityResult struct {
	Success      bool
	Damage       int
	StatusEffect *StatusEffectData
	Message      string
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
