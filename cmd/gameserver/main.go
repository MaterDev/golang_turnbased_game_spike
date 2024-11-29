package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/MaterDev/golang_turnbased_game_spike/internal/game"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
)

type BattleRequest struct {
	Character1 game.Character `json:"Character1"`
	Character2 game.Character `json:"Character2"`
}

// BattleResponse represents the JSON-safe version of a Battle
type BattleResponse struct {
	ID         string          `json:"ID"`
	Character1 *game.Character `json:"Character1"`
	Character2 *game.Character `json:"Character2"`
	State      game.BattleState `json:"State"`
	Winner     *game.Character `json:"Winner,omitempty"`
	Round      int            `json:"Round"`
}

// Convert Battle to BattleResponse
func toBattleResponse(b *game.Battle) BattleResponse {
	return BattleResponse{
		ID:         b.ID,
		Character1: b.Character1,
		Character2: b.Character2,
		State:      b.State,
		Winner:     b.Winner,
		Round:      b.Round,
	}
}

// BattleManager handles storing and retrieving battles
type BattleManager struct {
	battles map[string]*game.Battle
	mu      sync.RWMutex
}

func NewBattleManager() *BattleManager {
	return &BattleManager{
		battles: make(map[string]*game.Battle),
	}
}

func (bm *BattleManager) CreateBattle(char1, char2 *game.Character) *game.Battle {
	battle := game.NewBattle(char1, char2)
	bm.mu.Lock()
	bm.battles[battle.ID] = battle
	bm.mu.Unlock()
	return battle
}

func (bm *BattleManager) GetBattle(id string) *game.Battle {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return bm.battles[id]
}

var battleManager = NewBattleManager()

func main() {
	r := mux.NewRouter()

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Add logging middleware first to log all requests
	r.Use(loggingMiddleware)

	// Create API subrouter
	api := r.PathPrefix("/api").Subrouter()

	// Battle endpoints
	api.HandleFunc("/battles", createBattleHandler).Methods("POST", "OPTIONS")
	api.HandleFunc("/battles/{id}/start", startBattleHandler).Methods("POST", "OPTIONS")
	api.HandleFunc("/battles/{id}/action", submitActionHandler).Methods("POST", "OPTIONS")

	// Serve static files (for non-API routes)
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.HandleFunc("/", serveIndex)

	// Start server
	fmt.Println("Server starting on :3000...")
	log.Fatal(http.ListenAndServe(":3000", r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		log.Printf("Headers: %v", r.Header)
		
		// Read and log the request body for debugging
		if r.Method == "POST" {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading body: %v", err)
			} else {
				log.Printf("Request body: %s", string(bodyBytes))
				// Create a new reader with the same body content
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}
		
		next.ServeHTTP(w, r)
	})
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, filepath.Join("static", "index.html"))
}

func createBattleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Log request body for debugging
	var request BattleRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Initialize abilities if they're nil
	if request.Character1.Abilities == nil {
		log.Printf("Initializing abilities for Character1")
		request.Character1.Abilities = []game.Ability{
			{
				Name:        "Basic Attack",
				Damage:      10,
				CooldownMax: 0,
			},
			{
				Name:        "Power Strike",
				Damage:      20,
				CooldownMax: 2,
				StatusEffect: game.StatusEffectData{
					Type:     game.StatusEnraged,
					Duration: 2,
					Potency:  20,
				},
			},
		}
	}

	if request.Character2.Abilities == nil {
		log.Printf("Initializing abilities for Character2")
		request.Character2.Abilities = []game.Ability{
			{
				Name:        "Basic Attack",
				Damage:      8,
				CooldownMax: 0,
			},
			{
				Name:        "Fireball",
				Damage:      15,
				CooldownMax: 2,
				StatusEffect: game.StatusEffectData{
					Type:     game.StatusBurning,
					Duration: 3,
					Potency:  5,
				},
			},
		}
	}

	// Log the received characters
	log.Printf("Creating battle with characters: %+v vs %+v", request.Character1, request.Character2)

	// Set character IDs
	request.Character1.ID = uuid.New().String()
	request.Character2.ID = uuid.New().String()

	// Create new battle
	battle := battleManager.CreateBattle(&request.Character1, &request.Character2)
	if battle == nil {
		log.Printf("Error creating battle: battle is nil")
		http.Error(w, "Failed to create battle", http.StatusInternalServerError)
		return
	}

	// Log the created battle
	log.Printf("Created battle: %+v", battle)
	log.Printf("Character1 ID: %s", battle.Character1.ID)
	log.Printf("Character2 ID: %s", battle.Character2.ID)

	// Convert to response and return
	response := toBattleResponse(battle)
	log.Printf("Sending battle response: %+v", response)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func startBattleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	vars := mux.Vars(r)
	battleID := vars["id"]
	log.Printf("Starting battle with ID: %s", battleID)
	
	battle := battleManager.GetBattle(battleID)
	if battle == nil {
		http.Error(w, fmt.Sprintf("Battle not found: %s", battleID), http.StatusNotFound)
		return
	}

	if err := battle.Start(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to start battle: %v", err), http.StatusBadRequest)
		return
	}

	response := toBattleResponse(battle)
	log.Printf("Battle started successfully: %+v", response)
	json.NewEncoder(w).Encode(response)
}

func submitActionHandler(w http.ResponseWriter, r *http.Request) {
    // Check if this is an HTMX request
    if r.Header.Get("HX-Request") == "true" {
        w.Header().Set("Content-Type", "text/html")
    } else {
        w.Header().Set("Content-Type", "application/json")
    }
    
    vars := mux.Vars(r)
    battleID := vars["id"]
    
    battle := battleManager.GetBattle(battleID)
    if battle == nil {
        if r.Header.Get("HX-Request") == "true" {
            http.Error(w, fmt.Sprintf("Battle not found: %s", battleID), http.StatusNotFound)
        } else {
            http.Error(w, fmt.Sprintf("Battle not found: %s", battleID), http.StatusNotFound)
        }
        return
    }

    var action game.BattleAction
    if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
        if r.Header.Get("HX-Request") == "true" {
            http.Error(w, fmt.Sprintf("Invalid action: %v", err), http.StatusBadRequest)
        } else {
            http.Error(w, fmt.Sprintf("Invalid action: %v", err), http.StatusBadRequest)
        }
        return
    }

    // Handle target selection from form data for HTMX requests
    if r.Header.Get("HX-Request") == "true" {
        char1Target := r.FormValue("char1-target")
        char2Target := r.FormValue("char2-target")
        
        // Set target based on which character is acting
        if action.CharacterID == battle.Character1.ID {
            if char1Target == "self" {
                action.TargetID = battle.Character1.ID
            } else {
                action.TargetID = battle.Character2.ID
            }
        } else {
            if char2Target == "self" {
                action.TargetID = battle.Character2.ID
            } else {
                action.TargetID = battle.Character1.ID
            }
        }
    }

    result := battle.SubmitAction(action)
    
    if r.Header.Get("HX-Request") == "true" {
        // Return updated battle view HTML
        char1StatusEffects := ""
        if len(battle.Character1.StatusEffects) > 0 {
            effects := make([]string, len(battle.Character1.StatusEffects))
            for i, effect := range battle.Character1.StatusEffects {
                effects[i] = string(effect.Type)
            }
            char1StatusEffects = fmt.Sprintf("Status Effects: %s", strings.Join(effects, ", "))
        }

        char2StatusEffects := ""
        if len(battle.Character2.StatusEffects) > 0 {
            effects := make([]string, len(battle.Character2.StatusEffects))
            for i, effect := range battle.Character2.StatusEffects {
                effects[i] = string(effect.Type)
            }
            char2StatusEffects = fmt.Sprintf("Status Effects: %s", strings.Join(effects, ", "))
        }

        disabledAttr := ""
        if battle.State == game.BattleStateComplete {
            disabledAttr = "disabled"
        }

        tmpl := `
        <div id="battle-view">
            <div class="battle-container">
                <div class="character" id="char1">
                    <h2>%s</h2>
                    <div class="stats">
                        Health: %d<br>
                        Attack: %d<br>
                        Defense: %d<br>
                        Speed: %d
                    </div>
                    <div class="status-effects">%s</div>
                    <div class="actions">
                        <select id="char1-target" onchange="updateHtmxAttributes()">
                            <option value="enemy">Target Enemy</option>
                            <option value="self">Target Self</option>
                        </select>
                        <div class="ability-buttons">
                            <button class="basic-attack" 
                                    hx-post="/api/battles/%s/action"
                                    hx-trigger="click"
                                    hx-vals='{"CharacterID":"%s","AbilityIndex":0,"TargetID":"%s"}'
                                    hx-target="#battle-view"
                                    hx-swap="outerHTML"
                                    %s>
                                Basic Attack
                            </button>
                            <button class="special-attack"
                                    hx-post="/api/battles/%s/action"
                                    hx-trigger="click"
                                    hx-vals='{"CharacterID":"%s","AbilityIndex":1,"TargetID":"%s"}'
                                    hx-target="#battle-view"
                                    hx-swap="outerHTML"
                                    %s>
                                Power Strike
                            </button>
                        </div>
                    </div>
                </div>
                <div class="character" id="char2">
                    <h2>%s</h2>
                    <div class="stats">
                        Health: %d<br>
                        Attack: %d<br>
                        Defense: %d<br>
                        Speed: %d
                    </div>
                    <div class="status-effects">%s</div>
                    <div class="actions">
                        <select id="char2-target" onchange="updateHtmxAttributes()">
                            <option value="enemy">Target Enemy</option>
                            <option value="self">Target Self</option>
                        </select>
                        <div class="ability-buttons">
                            <button class="basic-attack"
                                    hx-post="/api/battles/%s/action"
                                    hx-trigger="click"
                                    hx-vals='{"CharacterID":"%s","AbilityIndex":0,"TargetID":"%s"}'
                                    hx-target="#battle-view"
                                    hx-swap="outerHTML"
                                    %s>
                                Basic Attack
                            </button>
                            <button class="special-attack"
                                    hx-post="/api/battles/%s/action"
                                    hx-trigger="click"
                                    hx-vals='{"CharacterID":"%s","AbilityIndex":1,"TargetID":"%s"}'
                                    hx-target="#battle-view"
                                    hx-swap="outerHTML"
                                    %s>
                                Fireball
                            </button>
                        </div>
                    </div>
                </div>
            </div>
            <div class="battle-log" id="battle-log">%s</div>
        </div>`

        // For character 1's buttons
        char1Target := battle.Character2.ID // Default to enemy
        if r.FormValue("char1-target") == "self" {
            char1Target = battle.Character1.ID
        }

        // For character 2's buttons
        char2Target := battle.Character1.ID // Default to enemy
        if r.FormValue("char2-target") == "self" {
            char2Target = battle.Character2.ID
        }

        battleLog := ""
        if result.Success {
            battleLog = fmt.Sprintf("<div>%s</div>", result.Message)
        }

        fmt.Fprintf(w, tmpl,
            // Character 1
            battle.Character1.Name,
            battle.Character1.Health,
            battle.Character1.Attack,
            battle.Character1.Defense,
            battle.Character1.Speed,
            char1StatusEffects,
            // Character 1 Basic Attack
            battle.ID, battle.Character1.ID, char1Target, disabledAttr,
            // Character 1 Special Attack
            battle.ID, battle.Character1.ID, char1Target, disabledAttr,
            // Character 2
            battle.Character2.Name,
            battle.Character2.Health,
            battle.Character2.Attack,
            battle.Character2.Defense,
            battle.Character2.Speed,
            char2StatusEffects,
            // Character 2 Basic Attack
            battle.ID, battle.Character2.ID, char2Target, disabledAttr,
            // Character 2 Special Attack
            battle.ID, battle.Character2.ID, char2Target, disabledAttr,
            // Battle log
            battleLog)
    } else {
        // Return JSON response for non-HTMX requests
        response := map[string]interface{}{
            "success": result.Success,
            "message": result.Message,
            "battle":  toBattleResponse(battle),
        }
        json.NewEncoder(w).Encode(response)
    }
}
