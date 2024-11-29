package game

import "testing"

func TestCharacter_ProcessStatusEffect(t *testing.T) {
	tests := []struct {
		name             string
		character        Character
		statusEffectData StatusEffectData
		expected         []int // Expected numerical values for stats after x rounds
		statToCheck     string // Which stat to check for the expected values
	}{
		{
			name: "accelerate effect",
			character: Character{
				Name:    "Speedster",
				Health:  100,
				Attack:  10,
				Defense: 10,
				Speed:   10,
			},
			statusEffectData: StatusEffectData{
				Type:     StatusAccelerate,
				Duration: 3,
				Potency:  10, // 10% increase per round
			},
			expected:    []int{11, 12, 13}, // Speed increases by 10% each round
			statToCheck: "speed",
		},
		{
			name: "burning effect",
			character: Character{
				Name:    "Tank",
				Health:  100,
				Attack:  10,
				Defense: 10,
				Speed:   10,
			},
			statusEffectData: StatusEffectData{
				Type:     StatusBurning,
				Duration: 3,
				Potency:  5, // 5% damage that increases with duration
			},
			expected:    []int{99, 99, 99}, // Damage increases each round but defense reduces it
			statToCheck: "health",
		},
		{
			name: "poison effect",
			character: Character{
				Name:    "Victim",
				Health:  100,
				Attack:  10,
				Defense: 10,
				Speed:   10,
			},
			statusEffectData: StatusEffectData{
				Type:     StatusPoisoned,
				Duration: 3,
				Potency:  15, // 15% damage that decreases with duration
			},
			expected:    []int{72, 41, 4}, // Damage decreases with duration
			statToCheck: "health",
		},
		{
			name: "enrage effect",
			character: Character{
				Name:    "Berserker",
				Health:  100,
				Attack:  10,
				Defense: 10,
				Speed:   10,
			},
			statusEffectData: StatusEffectData{
				Type:     StatusEnraged,
				Duration: 3,
				Potency:  20, // 20% attack increase
			},
			expected:    []int{12, 14, 16}, // Attack increases by 20% each round
			statToCheck: "attack",
		},
		{
			name: "regeneration effect",
			character: Character{
				Name:    "Healer",
				Health:  100,
				Attack:  10,
				Defense: 10,
				Speed:   10,
			},
			statusEffectData: StatusEffectData{
				Type:     StatusRegenerating,
				Duration: 3,
				Potency:  10, // 10% health recovery
			},
			expected:    []int{111, 123, 136}, // Health increases by 10% each round
			statToCheck: "health",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char := tt.character // Create a copy of the character
			char.StatusEffects = []StatusEffectData{tt.statusEffectData}

			results := []int{}
			for i := 0; i < tt.statusEffectData.Duration; i++ {
				char.ProcessStatusEffect()
				
				// Get the appropriate stat value based on the effect type
				var statValue int
				switch tt.statToCheck {
				case "health":
					statValue = char.Health
				case "attack":
					statValue = char.Attack
				case "defense":
					statValue = char.Defense
				case "speed":
					statValue = char.Speed
				}
				results = append(results, statValue)
			}

			// Compare results
			for i, expected := range tt.expected {
				if results[i] != expected {
					t.Errorf("Round %d: expected %d, got %d", i+1, expected, results[i])
				}
			}

			// Check that the effect was removed after duration expired
			if len(char.StatusEffects) != 0 {
				t.Error("Status effect should be removed after duration expires")
			}
		})
	}
}

func TestCharacter_MultipleStatusEffects(t *testing.T) {
	char := Character{
		Name:    "Multi-Effect Target",
		Health:  100,
		Attack:  10,
		Defense: 10,
		Speed:   10,
	}

	// Apply multiple status effects
	char.StatusEffects = []StatusEffectData{
		{
			Type:     StatusEnraged,
			Duration: 2,
			Potency:  20, // 20% attack increase
		},
		{
			Type:     StatusAccelerate,
			Duration: 2,
			Potency:  10, // 10% speed increase
		},
	}

	// Process effects and check results
	char.ProcessStatusEffect()

	// Check first round effects
	if char.Attack != 12 { // 10 + 20% increase
		t.Errorf("Expected attack to be 12, got %d", char.Attack)
	}
	if char.Speed != 11 { // 10 + 10% increase
		t.Errorf("Expected speed to be 11, got %d", char.Speed)
	}

	// Process second round
	char.ProcessStatusEffect()

	// Check second round effects
	if char.Attack != 14 { // Previous + 20% increase
		t.Errorf("Expected attack to be 14, got %d", char.Attack)
	}
	if char.Speed != 12 { // Previous + 10% increase
		t.Errorf("Expected speed to be 12, got %d", char.Speed)
	}

	// Check that both effects were removed
	if len(char.StatusEffects) != 0 {
		t.Errorf("Expected all status effects to be removed, got %d remaining", len(char.StatusEffects))
	}
}

func TestStatusEffect_DurationManagement(t *testing.T) {
	char := Character{
		Name:    "Duration Test",
		Health:  100,
		Attack:  10,
		Defense: 10,
		Speed:   10,
	}

	// Apply a status effect with 1 turn duration
	char.StatusEffects = []StatusEffectData{
		{
			Type:     StatusEnraged,
			Duration: 1,
			Potency:  20,
		},
	}

	// Process the effect
	char.ProcessStatusEffect()

	// Check that the effect was applied
	if char.Attack != 12 {
		t.Errorf("Expected attack to be 12, got %d", char.Attack)
	}

	// Check that the effect was removed after one turn
	if len(char.StatusEffects) != 0 {
		t.Error("Expected status effect to be removed after duration expired")
	}

	// Process again to ensure no further changes
	char.ProcessStatusEffect()
	if char.Attack != 12 {
		t.Errorf("Expected attack to remain 12, got %d", char.Attack)
	}
}
