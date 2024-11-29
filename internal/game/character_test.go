package game

import "testing"

func TestCharacter_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		character Character
		want      bool
	}{
		{
			name: "valid character",
			character: Character{
				ID:      "1",
				Name:    "Test",
				Health:  100,
				Attack:  10,
				Defense: 5,
				Speed:   10,
			},
			want: true,
		},
		{
			name: "invalid - empty name",
			character: Character{
				ID:      "2",
				Name:    "",
				Health:  100,
				Attack:  10,
				Defense: 5,
				Speed:   10,
			},
			want: false,
		},
		{
			name: "invalid - zero health",
			character: Character{
				ID:      "3",
				Name:    "Test",
				Health:  0,
				Attack:  10,
				Defense: 5,
				Speed:   10,
			},
			want: false,
		},
		{
			name: "invalid - negative attack",
			character: Character{
				ID:      "4",
				Name:    "Test",
				Health:  100,
				Attack:  -1,
				Defense: 5,
				Speed:   10,
			},
			want: false,
		},
		{
			name: "invalid - negative defense",
			character: Character{
				ID:      "5",
				Name:    "Test",
				Health:  100,
				Attack:  10,
				Defense: -1,
				Speed:   10,
			},
			want: false,
		},
		{
			name: "invalid - zero speed",
			character: Character{
				ID:      "6",
				Name:    "Test",
				Health:  100,
				Attack:  10,
				Defense: 5,
				Speed:   0,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.character.IsValid(); got != tt.want {
				t.Errorf("Character.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCharacter_TakeDamage(t *testing.T) {
	tests := []struct {
		name           string
		character      Character
		damage         int
		wantHealth     int
		wantMinHealth  bool
	}{
		{
			name: "normal damage",
			character: Character{
				Health:  100,
				Defense: 5,
			},
			damage:     20,
			wantHealth: 85, // 100 - (20 - 5)
		},
		{
			name: "damage fully blocked by defense",
			character: Character{
				Health:  100,
				Defense: 20,
			},
			damage:     15,
			wantHealth: 100,
		},
		{
			name: "lethal damage",
			character: Character{
				Health:  50,
				Defense: 0,
			},
			damage:        100,
			wantHealth:    0,
			wantMinHealth: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.character.TakeDamage(tt.damage)
			if tt.character.Health != tt.wantHealth {
				t.Errorf("Character.TakeDamage(%d) = %d health, want %d", 
					tt.damage, tt.character.Health, tt.wantHealth)
			}
			if tt.wantMinHealth && tt.character.Health < 0 {
				t.Error("Character health went below 0")
			}
		})
	}
}

func TestCharacter_UseAbility(t *testing.T) {
	tests := []struct {
		name           string
		attacker       Character
		target         Character
		abilityIndex   int
		wantSuccess    bool
		wantHealth     int
		wantEffect     bool
		wantCooldown   bool
	}{
		{
			name: "basic attack",
			attacker: Character{
				Name:    "Attacker",
				Attack:  20,
				Abilities: []Ability{
					{
						Name:        "Basic Attack",
						Damage:      10,
						CooldownMax: 0,
					},
				},
			},
			target: Character{
				Name:    "Target",
				Health:  100,
				Defense: 5,
			},
			abilityIndex: 0,
			wantSuccess:  true,
			wantHealth:   95,
			wantEffect:   false,
		},
		{
			name: "attack with status effect",
			attacker: Character{
				Name:    "Attacker",
				Attack:  20,
				Abilities: []Ability{
					{
						Name:        "Burning Strike",
						Damage:      15,
						CooldownMax: 2,
						StatusEffect: StatusEffectData{
							Type:     StatusBurning,
							Duration: 2,
							Potency:  10,
						},
					},
				},
			},
			target: Character{
				Name:    "Target",
				Health:  100,
				Defense: 5,
			},
			abilityIndex: 0,
			wantSuccess:  true,
			wantHealth:   90,
			wantEffect:   true,
			wantCooldown: true,
		},
		{
			name: "invalid ability index",
			attacker: Character{
				Name:    "Attacker",
				Attack:  20,
				Abilities: []Ability{
					{
						Name:   "Basic Attack",
						Damage: 10,
					},
				},
			},
			target: Character{
				Name:    "Target",
				Health:  100,
			},
			abilityIndex: 1,
			wantSuccess:  false,
		},
		{
			name: "ability on cooldown",
			attacker: Character{
				Name:    "Attacker",
				Attack:  20,
				Abilities: []Ability{
					{
						Name:        "Power Attack",
						Damage:      20,
						CooldownMax: 2,
						Cooldown:    2,
					},
				},
			},
			target: Character{
				Name:    "Target",
				Health:  100,
			},
			abilityIndex: 0,
			wantSuccess:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.attacker.UseAbility(tt.abilityIndex, &tt.target)
			if result.Success != tt.wantSuccess {
				t.Errorf("UseAbility() success = %v, want %v", result.Success, tt.wantSuccess)
			}

			if tt.wantSuccess {
				if tt.target.Health != tt.wantHealth {
					t.Errorf("Target health = %d, want %d", tt.target.Health, tt.wantHealth)
				}

				if tt.wantEffect && len(tt.target.StatusEffects) != 1 {
					t.Error("Expected status effect to be applied")
				}

				if tt.wantCooldown && tt.attacker.Abilities[tt.abilityIndex].Cooldown != tt.attacker.Abilities[tt.abilityIndex].CooldownMax {
					t.Error("Expected ability to be on cooldown")
				}
			}
		})
	}
}

func TestCharacter_GetEffectScalingValue(t *testing.T) {
	char := Character{
		Health:  100,
		Attack:  20,
		Defense: 10,
		Speed:   15,
	}

	tests := []struct {
		name      string
		statName  string
		scalar    int
		potency   int
		divisor   int
		want      int
	}{
		{
			name:     "scale health up",
			statName: "health",
			scalar:   100,
			potency:  20,
			divisor:  1,
			want:     120, // 100 + (100 * 20/100)
		},
		{
			name:     "scale attack up with divisor",
			statName: "attack",
			scalar:   100,
			potency:  50,
			divisor:  2,
			want:     15,  // (20 + (20 * 50/100)) / 2
		},
		{
			name:     "scale defense",
			statName: "defense",
			scalar:   100,
			potency:  10,
			divisor:  1,
			want:     11,  // 10 + (10 * 10/100)
		},
		{
			name:     "scale speed",
			statName: "speed",
			scalar:   100,
			potency:  30,
			divisor:  1,
			want:     19,  // 15 + (15 * 30/100)
		},
		{
			name:     "unknown stat",
			statName: "unknown",
			scalar:   100,
			potency:  10,
			divisor:  1,
			want:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := char.GetEffectScalingValue(tt.statName, tt.scalar, tt.potency, tt.divisor); got != tt.want {
				t.Errorf("GetEffectScalingValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
