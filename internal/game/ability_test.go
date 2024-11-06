package game

import "testing"

// Testing that ability cant be used when cooldown is greater tahn 0
func TestAbility_CanUse(t *testing.T) {
	ability := Ability{
		Name:        "Fireball",
		Damage:      30,
		CooldownMax: 3,
		Cooldown:    0,
		StatusEffect: StatusEffectData{
			Type:     StatusBurning,
			Duration: 2,
			Potency:  5,
		},
	}

	// Test initial use (if ability can be used)
	if !ability.CanUse() {
		t.Error("Expected ability to be usable when cooldown is 0")
	}

	// Use the ability and test the cooldown
	ability.Use()
	if ability.CanUse() {
		t.Error("Expected ability to be unusable right after use")
	}

	// Test cooldown reduction
	ability.ReduceCooldown()
	if ability.Cooldown != 2 {
		t.Errorf("Expected cooldown to be 2, got %d", ability.Cooldown)
	}
}

// Testing damage of ability being used & status effect.
func TestCharacter_UseAbility(t *testing.T) {
	attacker := Character{
		Name:    "Mage",
		Health:  80,
		Attack:  20,
		Defense: 5,
		Speed:   8,
		Abilities: []Ability{
			{
				Name:        "Fireball",
				Damage:      30,
				CooldownMax: 3,
				Cooldown:    0,
				StatusEffect: StatusEffectData{
					Type:     StatusBurning,
					Duration: 2,
					Potency:  5,
				},
			},
		},
	}

	target := Character{
		Name:      "Punchin' Bag",
		Health:    1000,
		Attack:    0,
		Defense:   0,
		Speed:     0,
		Abilities: []Ability{},
	}

	result := attacker.UseAbility(0, &target)
	if !result.Success {
		t.Error("Expected ability to be used successfully!")
	}

	// Check damage calculation

	// Check status effect application
}
