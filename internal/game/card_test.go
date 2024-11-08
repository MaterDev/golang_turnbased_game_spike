package game

import "testing"

func TestCharacter_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		character Character
		expected  bool
	}{
		{
			name: "valid Primal Shifter character",
			character: Character{
				ID:        1,
				Name:      "Primal Shifter",
				Health:    100,
				Attack:    15,
				Defense:   10,
				Speed:     7,
				Abilities: []Ability{},
			},
			expected: true,
		},
		{
			name: "invalid character - negative health",
			character: Character{
				ID:        2,
				Name:      "Invalid",
				Health:    -10,
				Attack:    15,
				Defense:   10,
				Speed:     7,
				Abilities: []Ability{},
			},
			expected: false,
		},
		{
			name: "invalid character - empty name",
			character: Character{
				ID:        1,
				Name:      "",
				Health:    100,
				Attack:    15,
				Defense:   10,
				Speed:     7,
				Abilities: []Ability{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.character.IsValid(); got != tt.expected {
				t.Errorf("Character.IsValid() = %v, wanted %v", got, tt.expected)
			}
		})
	}
}

func TestCharacter_TakeDamage(t *testing.T) {
	character := Character{
		Name:          "Primal Shifter",
		Abilities:     []Ability{},
		StatusEffects: []StatusEffectData{},
		Health:        100,
		Attack:        15,
		Defense:       10,
		Speed:         7,
	}

	// Test basic damage calculation
	damage := 20
	character.TakeDamage(damage)
	expectedHealth := 90 // 100 - (20- 10 defense)

	if character.Health != expectedHealth {
		t.Errorf("Character.TakeDamage(%d) = %d health, want %d", damage, character.Health, expectedHealth)
	}
}
