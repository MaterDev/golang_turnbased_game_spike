import { useState, useCallback } from 'react';
import { Battle, Character } from './types';
import { CharacterForm } from './components/CharacterForm';
import { BattleView } from './components/BattleView';
import { API_BASE_URL } from './config';
import './App.css';

const defaultCharacter1: Character = {
  ID: '',
  Name: 'Warrior',
  Health: 100,
  Attack: 15,
  Defense: 10,
  Speed: 8,
  StatusEffects: [],
  Abilities: [
    {
      Name: 'Basic Attack',
      Damage: 10,
      CooldownMax: 0,
    },
    {
      Name: 'Power Strike',
      Damage: 20,
      CooldownMax: 2,
      StatusEffect: {
        Type: 'ENRAGED',
        Duration: 2,
        Potency: 20,
      },
    },
  ],
};

const defaultCharacter2: Character = {
  ID: '',
  Name: 'Mage',
  Health: 80,
  Attack: 20,
  Defense: 5,
  Speed: 12,
  StatusEffects: [],
  Abilities: [
    {
      Name: 'Basic Attack',
      Damage: 8,
      CooldownMax: 0,
    },
    {
      Name: 'Fireball',
      Damage: 15,
      CooldownMax: 2,
      StatusEffect: {
        Type: 'BURNING',
        Duration: 3,
        Potency: 5,
      },
    },
  ],
};

function App() {
  const [char1, setChar1] = useState<Character>(defaultCharacter1);
  const [char2, setChar2] = useState<Character>(defaultCharacter2);
  const [battle, setBattle] = useState<Battle | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const resetBattle = useCallback(() => {
    setBattle(null);
    setError(null);
    setIsLoading(false);
  }, []);

  const createBattle = async () => {
    setIsLoading(true);
    setError(null);
    try {
      console.log('Creating battle with characters:', { char1, char2 });
      const response = await fetch(`${API_BASE_URL}/battles`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          Character1: char1,
          Character2: char2,
        }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(`Failed to create battle: ${errorText}`);
      }

      const newBattle = await response.json();
      console.log('Battle created:', newBattle);

      // Start the battle
      const startResponse = await fetch(`${API_BASE_URL}/battles/${newBattle.ID}/start`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!startResponse.ok) {
        const errorText = await startResponse.text();
        throw new Error(`Failed to start battle: ${errorText}`);
      }

      const startedBattle = await startResponse.json();
      console.log('Battle started:', startedBattle);
      setBattle(startedBattle);
      setError(null);
    } catch (err) {
      console.error('Battle creation error:', err);
      setError(err instanceof Error ? err.message : 'An error occurred');
      setBattle(null);
    } finally {
      setIsLoading(false);
    }
  };

  const submitAction = async (action: { CharacterID: string; AbilityIndex: number; TargetID: string }) => {
    if (!battle?.ID) {
      console.error('Cannot submit action: No active battle');
      return;
    }

    console.log('Starting action submission:', {
      battleId: battle.ID,
      action,
      currentBattleState: battle.State,
      currentRound: battle.Round
    });

    setIsLoading(true);
    try {
      const response = await fetch(`${API_BASE_URL}/battles/${battle.ID}/action`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(action),
      });

      console.log('Action API response status:', response.status);

      if (!response.ok) {
        const errorText = await response.text();
        console.error('API error response:', errorText);
        throw new Error(`Failed to submit action: ${errorText}`);
      }

      const actionResponse = await response.json();
      console.log('Raw API response:', actionResponse);

      // Check if the action was successful
      if (!actionResponse.success) {
        throw new Error(actionResponse.message || 'Action failed');
      }

      // If the action response includes the updated battle state, use it
      if (actionResponse.battle && actionResponse.battle.ID === battle.ID) {
        console.log('Using battle state from action response');
        setBattle(actionResponse.battle);
        setError(null);
      } else {
        // Otherwise fetch the latest battle state
        console.log('Fetching fresh battle state');
        const battleResponse = await fetch(`${API_BASE_URL}/battles/${battle.ID}`, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        });

        if (!battleResponse.ok) {
          // If we can't get the battle state, but the action was successful,
          // we'll keep the current battle state and show a warning
          console.warn('Could not fetch updated battle state, keeping current state');
          setError('Battle state may be outdated - please refresh');
          return;
        }

        const updatedBattle = await battleResponse.json();
        console.log('Battle state updated:', {
          previousState: battle.State,
          newState: updatedBattle.State,
          previousRound: battle.Round,
          newRound: updatedBattle.Round,
          char1: updatedBattle.Character1,
          char2: updatedBattle.Character2,
          winner: updatedBattle.Winner?.Name
        });
        
        setBattle(updatedBattle);
        setError(null);
      }
    } catch (err) {
      console.error('Action submission error:', err);
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="App">
      <h1>Turn-Based Battle</h1>
      
      {error && (
        <div className="error">
          {error}
          <button onClick={resetBattle} className="reset-button">
            Start New Battle
          </button>
        </div>
      )}
      
      {isLoading && <div className="loading">Loading...</div>}

      {!battle ? (
        <>
          <h2>Create New Battle</h2>
          <div className="battle-container">
            <CharacterForm
              label="Character 1"
              defaultCharacter={defaultCharacter1}
              onChange={setChar1}
            />
            <CharacterForm
              label="Character 2"
              defaultCharacter={defaultCharacter2}
              onChange={setChar2}
            />
          </div>
          <button 
            onClick={createBattle} 
            disabled={isLoading}
            className="create-battle-button"
          >
            {isLoading ? 'Creating Battle...' : 'Create Battle'}
          </button>
        </>
      ) : (
        <>
          <BattleView 
            battle={battle} 
            onAction={submitAction} 
          />
          {battle.State === 'COMPLETE' && (
            <button onClick={resetBattle} className="reset-button">
              Start New Battle
            </button>
          )}
        </>
      )}
    </div>
  );
}

export default App;
