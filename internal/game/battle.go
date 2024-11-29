package game

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type BattleState string

const (
	BattleStatePending  BattleState = "PENDING"
	BattleStateActive   BattleState = "ACTIVE"
	BattleStateComplete BattleState = "COMPLETE"
)

type Battle struct {
	ID         string
	Character1 *Character
	Character2 *Character
	State      BattleState
	Winner     *Character
	Round      int
	mu         sync.Mutex
	ActionChan chan BattleAction
}

type BattleAction struct {
	CharacterID   string
	AbilityIndex  int
	TargetID      string
	ResponseChan  chan BattleActionResult
}

type BattleActionResult struct {
	Success bool
	Message string
	Battle  *Battle
}

func NewBattle(char1, char2 *Character) *Battle {
	return &Battle{
		ID:         uuid.New().String(),
		Character1: char1,
		Character2: char2,
		State:      BattleStatePending,
		Round:      1,
		ActionChan: make(chan BattleAction, 100),
	}
}

func (b *Battle) Start() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.State != BattleStatePending {
		return errors.New("battle already started")
	}

	b.State = BattleStateActive
	
	// Start battle loop in goroutine
	go b.battleLoop()
	
	return nil
}

func (b *Battle) battleLoop() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case action := <-b.ActionChan:
			result := b.processAction(action)
			if action.ResponseChan != nil {
				action.ResponseChan <- result
			}
			
			// Process status effects at the end of each round
			if b.State == BattleStateActive {
				b.Character1.ProcessStatusEffect()
				b.Character2.ProcessStatusEffect()
			}

		case <-ticker.C:
			// Check battle state
			if b.State == BattleStateComplete {
				return
			}
		}
	}
}

func (b *Battle) processAction(action BattleAction) BattleActionResult {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.State != BattleStateActive {
		return BattleActionResult{
			Success: false,
			Message: "battle not active",
			Battle:  b,
		}
	}

	// Get acting character
	var actor *Character
	if b.Character1.ID == action.CharacterID {
		actor = b.Character1
	} else if b.Character2.ID == action.CharacterID {
		actor = b.Character2
	} else {
		return BattleActionResult{
			Success: false,
			Message: "invalid character ID",
			Battle:  b,
		}
	}

	// Get target character based on TargetID
	var target *Character
	if b.Character1.ID == action.TargetID {
		target = b.Character1
	} else if b.Character2.ID == action.TargetID {
		target = b.Character2
	} else {
		return BattleActionResult{
			Success: false,
			Message: "invalid target ID",
			Battle:  b,
		}
	}

	// Process the ability
	result := actor.UseAbility(action.AbilityIndex, target)
	if !result.Success {
		return BattleActionResult{
			Success: false,
			Message: result.Message,
			Battle:  b,
		}
	}

	// Check for battle end
	if b.Character1.Health <= 0 {
		b.Winner = b.Character2
		b.State = BattleStateComplete
	} else if b.Character2.Health <= 0 {
		b.Winner = b.Character1
		b.State = BattleStateComplete
	}

	return BattleActionResult{
		Success: true,
		Message: result.Message,
		Battle:  b,
	}
}

func (b *Battle) SubmitAction(action BattleAction) BattleActionResult {
	responseChan := make(chan BattleActionResult)
	action.ResponseChan = responseChan
	b.ActionChan <- action
	return <-responseChan
}
