package game

import (
	"testing"
	"time"
)

func createTestCharacter(name string, health int) *Character {
	return &Character{
		ID:      name + "_id",
		Name:    name,
		Health:  health,
		Attack:  10,
		Defense: 5,
		Speed:   10,
		Abilities: []Ability{
			{
				Name:        "Basic Attack",
				Damage:      10,
				CooldownMax: 0,
			},
			{
				Name:        "Special Attack",
				Damage:      20,
				CooldownMax: 2,
				StatusEffect: StatusEffectData{
					Type:     StatusBurning,
					Duration: 2,
					Potency:  10,
				},
			},
		},
	}
}

func TestNewBattle(t *testing.T) {
	char1 := createTestCharacter("Warrior", 100)
	char2 := createTestCharacter("Mage", 80)

	battle := NewBattle(char1, char2)

	if battle.ID == "" {
		t.Error("Expected battle ID to be generated")
	}
	if battle.State != BattleStatePending {
		t.Errorf("Expected initial state to be PENDING, got %v", battle.State)
	}
	if battle.Round != 1 {
		t.Errorf("Expected initial round to be 1, got %d", battle.Round)
	}
	if battle.ActionChan == nil {
		t.Error("Expected ActionChan to be initialized")
	}
}

func TestBattle_Start(t *testing.T) {
	char1 := createTestCharacter("Warrior", 100)
	char2 := createTestCharacter("Mage", 80)
	battle := NewBattle(char1, char2)

	if err := battle.Start(); err != nil {
		t.Errorf("Unexpected error starting battle: %v", err)
	}

	if battle.State != BattleStateActive {
		t.Errorf("Expected battle state to be ACTIVE, got %v", battle.State)
	}

	// Test starting an already started battle
	if err := battle.Start(); err == nil {
		t.Error("Expected error when starting an already started battle")
	}
}

func TestBattle_SubmitAction(t *testing.T) {
	char1 := createTestCharacter("Warrior", 100)
	char2 := createTestCharacter("Mage", 80)
	battle := NewBattle(char1, char2)

	// Start the battle
	if err := battle.Start(); err != nil {
		t.Fatalf("Failed to start battle: %v", err)
	}

	tests := []struct {
		name        string
		action      BattleAction
		wantSuccess bool
	}{
		{
			name: "valid basic attack",
			action: BattleAction{
				CharacterID:  char1.ID,
				AbilityIndex: 0,
				TargetID:     char2.ID,
			},
			wantSuccess: true,
		},
		{
			name: "invalid character ID",
			action: BattleAction{
				CharacterID:  "invalid_id",
				AbilityIndex: 0,
				TargetID:     char2.ID,
			},
			wantSuccess: false,
		},
		{
			name: "invalid ability index",
			action: BattleAction{
				CharacterID:  char1.ID,
				AbilityIndex: 99,
				TargetID:     char2.ID,
			},
			wantSuccess: false,
		},
		{
			name: "invalid target ID",
			action: BattleAction{
				CharacterID:  char1.ID,
				AbilityIndex: 0,
				TargetID:     "invalid_target",
			},
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := battle.SubmitAction(tt.action)
			if result.Success != tt.wantSuccess {
				t.Errorf("SubmitAction() success = %v, want %v", result.Success, tt.wantSuccess)
			}
		})
	}
}

func TestBattle_BattleEnd(t *testing.T) {
	// Create a character with very low health and no defense
	char1 := createTestCharacter("Warrior", 10)
	char1.Defense = 0 // Remove defense to ensure lethal damage

	char2 := createTestCharacter("Mage", 80)
	char2.Attack = 50 // Increase attack to ensure lethal damage
	battle := NewBattle(char1, char2)

	// Start the battle
	if err := battle.Start(); err != nil {
		t.Fatalf("Failed to start battle: %v", err)
	}

	// Use special attack that should deal lethal damage
	result := battle.SubmitAction(BattleAction{
		CharacterID:  char2.ID,
		AbilityIndex: 1, // Special Attack
		TargetID:     char1.ID,
	})

	if !result.Success {
		t.Fatalf("Attack failed: %v", result.Message)
	}

	// Wait a bit for the battle loop to process
	time.Sleep(100 * time.Millisecond)

	if battle.State != BattleStateComplete {
		t.Errorf("Expected battle to be complete after lethal damage, got state: %v", battle.State)
	}

	if battle.Winner == nil {
		t.Error("Expected winner to be set")
	} else if battle.Winner.ID != char2.ID {
		t.Errorf("Expected winner to be char2, got %v", battle.Winner.Name)
	}
}

func TestBattle_ConcurrentActions(t *testing.T) {
	char1 := createTestCharacter("Warrior", 100)
	char2 := createTestCharacter("Mage", 100)
	battle := NewBattle(char1, char2)

	if err := battle.Start(); err != nil {
		t.Fatalf("Failed to start battle: %v", err)
	}

	// Submit multiple actions concurrently
	done := make(chan bool)
	go func() {
		for i := 0; i < 5; i++ {
			battle.SubmitAction(BattleAction{
				CharacterID:  char1.ID,
				AbilityIndex: 0,
				TargetID:     char2.ID,
			})
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 5; i++ {
			battle.SubmitAction(BattleAction{
				CharacterID:  char2.ID,
				AbilityIndex: 0,
				TargetID:     char1.ID,
			})
		}
		done <- true
	}()

	// Wait for both goroutines to finish
	<-done
	<-done

	// Verify both characters took damage
	if char1.Health >= 100 {
		t.Error("Expected char1 to take damage")
	}
	if char2.Health >= 100 {
		t.Error("Expected char2 to take damage")
	}
}
