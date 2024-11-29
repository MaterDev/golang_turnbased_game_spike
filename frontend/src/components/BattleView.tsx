import { useState, useCallback } from 'react';
import { Battle, BattleAction, Character } from '../types';

type BattleViewProps = {
    battle: Battle;
    onAction: (action: BattleAction) => Promise<void>;
};

export type CharacterViewProps = {
    character: Character;
    opponent: Character;
    onAction: (action: BattleAction) => Promise<void>;
    disabled: boolean;
    isCurrentTurn: boolean;
};

export function CharacterView({ character, opponent, onAction, disabled, isCurrentTurn }: CharacterViewProps) {
    const [targetSelf, setTargetSelf] = useState(false);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const targetId = targetSelf ? character?.ID : opponent?.ID;

    const handleAction = useCallback(async (abilityIndex: number) => {
        if (!character?.ID || !targetId || isSubmitting) {
            console.log('Action blocked:', { 
                reason: !character?.ID ? 'No character ID' : !targetId ? 'No target ID' : 'Already submitting',
                characterId: character?.ID,
                targetId,
                isSubmitting
            });
            return;
        }
        
        console.log('Attempting action:', {
            character: character.Name,
            ability: character.Abilities?.[abilityIndex]?.Name,
            target: targetSelf ? 'self' : 'opponent',
            abilityIndex
        });
        
        setIsSubmitting(true);
        try {
            await onAction({
                CharacterID: character.ID,
                AbilityIndex: abilityIndex,
                TargetID: targetId
            });
            console.log('Action completed successfully');
        } catch (error) {
            console.error('Action failed:', error);
        } finally {
            setIsSubmitting(false);
        }
    }, [character?.ID, targetId, onAction, isSubmitting, character, targetSelf]);

    if (!character || !opponent) {
        return <div>Loading character data...</div>;
    }

    return (
        <div className={`character-card ${isCurrentTurn ? 'current-turn' : ''}`}>
            <h2>{character.Name}</h2>
            <div className="stats">
                <div>Health: {character.Health}</div>
                <div>Attack: {character.Attack}</div>
                <div>Defense: {character.Defense}</div>
                <div>Speed: {character.Speed}</div>
            </div>
            <div className="status-effects">
                {character.StatusEffects?.length > 0 && (
                    <>Status Effects: {character.StatusEffects.map(e => e.Type).join(', ')}</>
                )}
            </div>
            <div className="actions">
                <select
                    value={targetSelf ? 'self' : 'enemy'}
                    onChange={(e) => setTargetSelf(e.target.value === 'self')}
                    disabled={disabled || !isCurrentTurn || isSubmitting}
                >
                    <option value="enemy">Target Enemy</option>
                    <option value="self">Target Self</option>
                </select>
                <div className="ability-buttons">
                    <button
                        className="basic-attack"
                        onClick={() => handleAction(0)}
                        disabled={disabled || !isCurrentTurn || isSubmitting || !targetId}
                        data-testid="basic-attack"
                    >
                        Basic Attack
                    </button>
                    <button
                        className="special-attack"
                        onClick={() => handleAction(1)}
                        disabled={disabled || !isCurrentTurn || isSubmitting || !targetId}
                        data-testid="special-attack"
                    >
                        {character.Abilities?.[1]?.Name || 'Special Attack'}
                    </button>
                </div>
            </div>
        </div>
    );
}

export function BattleView({ battle, onAction }: BattleViewProps) {
    const [lastActionTime, setLastActionTime] = useState(0);

    const handleAction = useCallback(async (action: BattleAction) => {
        console.log('BattleView handling action:', {
            action,
            currentState: battle.State,
            round: battle.Round
        });
        
        try {
            await onAction(action);
            setLastActionTime(prev => prev + 1);
            console.log('Action processed, new lastActionTime:', lastActionTime + 1);
        } catch (error) {
            console.error('Error performing action:', error);
        }
    }, [onAction, battle.State, battle.Round, lastActionTime]);

    if (!battle?.Character1 || !battle?.Character2) {
        console.log('Battle data missing:', { 
            char1Present: !!battle?.Character1, 
            char2Present: !!battle?.Character2 
        });
        return <div>Loading battle data...</div>;
    }

    // Determine whose turn it is based on speed
    const char1Speed = battle.Character1.Speed ?? 0;
    const char2Speed = battle.Character2.Speed ?? 0;
    
    // If Character2 has higher speed, they go first (lastActionTime === 0)
    // Then alternate turns based on lastActionTime
    const isChar2Turn = char2Speed > char1Speed ? (lastActionTime % 2 === 0) : (lastActionTime % 2 === 1);
    const isChar1Turn = !isChar2Turn;

    console.log('Turn state:', {
        lastActionTime,
        char1Speed,
        char2Speed,
        isChar1Turn,
        isChar2Turn,
        battleState: battle.State
    });

    const disabled = battle.State !== 'ACTIVE';

    return (
        <div id="battle-view">
            <div className="battle-status">
                <div>Round {battle.Round} - {battle.State}</div>
                {battle.State === 'COMPLETE' && battle.Winner && (
                    <div className="winner">Winner: {battle.Winner.Name}</div>
                )}
            </div>
            <div className="battle-container">
                <CharacterView
                    character={battle.Character1}
                    opponent={battle.Character2}
                    onAction={handleAction}
                    disabled={disabled}
                    isCurrentTurn={battle.State === 'ACTIVE' && isChar1Turn}
                />
                <CharacterView
                    character={battle.Character2}
                    opponent={battle.Character1}
                    onAction={handleAction}
                    disabled={disabled}
                    isCurrentTurn={battle.State === 'ACTIVE' && isChar2Turn}
                />
            </div>
            {battle.State === 'COMPLETE' && battle.Winner && (
                <div className="battle-log">
                    Battle ended! Winner: {battle.Winner.Name}
                </div>
            )}
        </div>
    );
}
