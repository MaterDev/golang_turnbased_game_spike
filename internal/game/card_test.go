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
				ID:          1,
				Name:        "Primal Shifter",
				Health:      100,
				Attack:      15,
				Defence:     10,
				Speed:       7,
				SpecialMove: "Whirlwind Strike",
			},
			expected: true,
		},
		{
			name: "invalid character - negative health",
			character: Character{
				ID:          2,
				Name:        "Invalid",
				Health:      -10,
				Attack:      15,
				Defence:     10,
				Speed:       7,
				SpecialMove: "None",
			},
			expected: false,
		},
		{
			name: "invalid character - empty name",
			character: Character{
				ID:          1,
				Name:        "",
				Health:      100,
				Attack:      15,
				Defence:     10,
				Speed:       7,
				SpecialMove: "Fireball",
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
