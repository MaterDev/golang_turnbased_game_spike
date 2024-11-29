package game

import "testing"

func TestAbility_CanUse(t *testing.T) {
	tests := []struct {
		name     string
		ability  Ability
		want     bool
	}{
		{
			name: "no cooldown ability",
			ability: Ability{
				Name:        "Basic Attack",
				Damage:      10,
				CooldownMax: 0,
				Cooldown:    0,
			},
			want: true,
		},
		{
			name: "ability ready",
			ability: Ability{
				Name:        "Special Attack",
				Damage:      20,
				CooldownMax: 2,
				Cooldown:    0,
			},
			want: true,
		},
		{
			name: "ability on cooldown",
			ability: Ability{
				Name:        "Power Move",
				Damage:      30,
				CooldownMax: 3,
				Cooldown:    2,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ability.CanUse(); got != tt.want {
				t.Errorf("Ability.CanUse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAbility_Use(t *testing.T) {
	tests := []struct {
		name         string
		ability      Ability
		wantCooldown int
	}{
		{
			name: "use basic attack",
			ability: Ability{
				Name:        "Basic Attack",
				Damage:      10,
				CooldownMax: 0,
			},
			wantCooldown: 0,
		},
		{
			name: "use special ability",
			ability: Ability{
				Name:        "Special Attack",
				Damage:      20,
				CooldownMax: 2,
			},
			wantCooldown: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ability.Use()
			if tt.ability.Cooldown != tt.wantCooldown {
				t.Errorf("After Use(), Cooldown = %v, want %v", tt.ability.Cooldown, tt.wantCooldown)
			}
		})
	}
}

func TestAbility_ReduceCooldown(t *testing.T) {
	tests := []struct {
		name         string
		ability      Ability
		reductions   int
		wantCooldown int
	}{
		{
			name: "reduce from max",
			ability: Ability{
				Name:        "Power Move",
				CooldownMax: 3,
				Cooldown:    3,
			},
			reductions:   2,
			wantCooldown: 1,
		},
		{
			name: "reduce to zero",
			ability: Ability{
				Name:        "Special Attack",
				CooldownMax: 2,
				Cooldown:    1,
			},
			reductions:   1,
			wantCooldown: 0,
		},
		{
			name: "already at zero",
			ability: Ability{
				Name:        "Basic Attack",
				CooldownMax: 1,
				Cooldown:    0,
			},
			reductions:   1,
			wantCooldown: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.reductions; i++ {
				tt.ability.ReduceCooldown()
			}
			if tt.ability.Cooldown != tt.wantCooldown {
				t.Errorf("After %d reductions, Cooldown = %v, want %v", 
					tt.reductions, tt.ability.Cooldown, tt.wantCooldown)
			}
		})
	}
}

// Testing status effect processing
func Test_ProcessStatusEffect(t *testing.T) {
	target := Character{
		Name:    "Mage",
		Health:  100,
		Attack:  10,
		Defense: 10,
		Speed:   10,
	}

	tests := []struct {
		name             string
		statusEffectData StatusEffectData
		expected         []int // Expected numerical values for stats after x rounds
		statToCheck     string // Which stat to check for the expected values
	}{
		{
			name: "Accelerate status applied - 3X",
			statusEffectData: StatusEffectData{
				Type:     StatusAccelerate,
				Duration: 3,
				Potency:  10, // 10% increase per round
			},
			expected:    []int{11, 12, 13}, // Speed increases by 10% each round
			statToCheck: "speed",
		},
		{
			name: "Burning status applied - 3X",
			statusEffectData: StatusEffectData{
				Type:     StatusBurning,
				Duration: 3,
				Potency:  5, // 5% damage that increases with duration
			},
			expected:    []int{99, 99, 99}, // Damage increases each round but defense reduces it
			statToCheck: "health",
		},
		{
			name: "Poison status applied - 3X",
			statusEffectData: StatusEffectData{
				Type:     StatusPoisoned,
				Duration: 3,
				Potency:  15, // 15% damage that decreases with duration
			},
			expected:    []int{72, 41, 4}, // Damage decreases with duration
			statToCheck: "health",
		},
		{
			name: "Enraged status applied - 3X",
			statusEffectData: StatusEffectData{
				Type:     StatusEnraged,
				Duration: 3,
				Potency:  20, // 20% attack increase
			},
			expected:    []int{12, 14, 16}, // Attack increases by 20% each round
			statToCheck: "attack",
		},
		{
			name: "Recovering status applied - 3X",
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
			// Reset target to initial state
			target = Character{
				Name:    "Mage",
				Health:  100,
				Attack:  10,
				Defense: 10,
				Speed:   10,
			}
			
			// Apply the status effect
			target.StatusEffects = []StatusEffectData{tt.statusEffectData}

			results := []int{}
			for i := 0; i < tt.statusEffectData.Duration; i++ {
				target.ProcessStatusEffect()
				
				// Get the appropriate stat value based on the effect type
				var statValue int
				switch tt.statToCheck {
				case "health":
					statValue = target.Health
				case "attack":
					statValue = target.Attack
				case "defense":
					statValue = target.Defense
				case "speed":
					statValue = target.Speed
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
			if len(target.StatusEffects) != 0 {
				t.Error("Status effect should be removed after duration expires")
			}
		})
	}
}
