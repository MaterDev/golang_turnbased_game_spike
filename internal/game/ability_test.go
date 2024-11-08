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
					Duration: 1,
					Potency:  10,
				},
			},
		},
	}

	target := Character{
		Name:    "Punchin' Bag",
		Health:  1000,
		Attack:  0,
		Defense: 0,
		Speed:   0,
	}

	result := attacker.UseAbility(0, &target)
	if !result.Success {
		t.Error("Expected ability to be used successfully!")
	}

	// Check damage calculation
	expectedHealth := 970 // 1000 - 30
	if target.Health != expectedHealth {
		t.Errorf("Expected target health to be %d, got %d", expectedHealth, target.Health)
	}

	// Check status effect application
	if len(target.StatusEffects) != 1 {
		t.Error("Expected status effect to be applied")
	}
}

// Testing status effect processing
func Test_ProcessStatusEffect(t *testing.T) {
	target := Character{
		Name:    "Mage",
		Health:  10,
		Attack:  10,
		Defense: 10,
		Speed:   10,
	}

	tests := []struct {
		name             string
		statusEffectData StatusEffectData
		expected         []int // Expected numerical values for stats after x rounds
	}{
		{
			name: "Accelerate status applied - 3X",
			statusEffectData: StatusEffectData{
				Type:     StatusAccelerate,
				Duration: 3,
				Potency:  10,
			},
			expected: []int{13, 16, 20},
		},
		{
			name: "Burning status applied - 3X",
			statusEffectData: StatusEffectData{
				Type:     StatusBurning,
				Duration: 3,
				Potency:  10,
			},
			expected: []int{13, 16, 20},
		},
		{
			name: "Poison status applied - 3X",
			statusEffectData: StatusEffectData{
				Type:     StatusPoisoned,
				Duration: 3,
				Potency:  10,
			},
			expected: []int{13, 16, 20},
		},
		{
			name: "Enraged status applied - 3X",
			statusEffectData: StatusEffectData{
				Type:     StatusEnraged,
				Duration: 3,
				Potency:  10,
			},
			expected: []int{13, 16, 20},
		},
		{
			name: "Recovering status applied - 3X",
			statusEffectData: StatusEffectData{
				Type:     StatusRecovering,
				Duration: 3,
				Potency:  10,
			},
			expected: []int{13, 16, 20},
		},
	}

	// TODO: Add test for compound status effect processing, since more than one can be applied at the same time.

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize empty []int{}
			// change character stateStatusEffects to be statusEffectData
			// Loop over range of duration
			// Call character.Test_ProcessStatusEffect
			// Append stat-after-change to statechangeArray
			// Compare to Expected
		})
	}
}
