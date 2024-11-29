import { useState } from 'react';
import { Character } from '../types';

type CharacterFormProps = {
    label: string;
    defaultCharacter: Character;
    onChange: (character: Character) => void;
};

export function CharacterForm({ label, defaultCharacter, onChange }: CharacterFormProps) {
    const [character, setCharacter] = useState(defaultCharacter);

    const handleChange = (field: keyof Character, value: string | number) => {
        let finalValue = value;
        
        // Clamp numeric values
        if (typeof value === 'number') {
            if (field === 'Health') {
                finalValue = Math.max(1, Math.min(999, value));
            } else {
                finalValue = Math.max(1, Math.min(100, value));
            }
        }

        const newCharacter = { ...character, [field]: finalValue };
        setCharacter(newCharacter);
        onChange(newCharacter);
    };

    return (
        <div className="character-form">
            <h3>{label}</h3>
            <div className="form-group">
                <label htmlFor={`${label}-name`}>Name</label>
                <input
                    id={`${label}-name`}
                    type="text"
                    value={character.Name}
                    onChange={(e) => handleChange('Name', e.target.value)}
                    placeholder="Enter character name"
                />
            </div>
            <div className="form-group">
                <label htmlFor={`${label}-health`}>Health</label>
                <input
                    id={`${label}-health`}
                    type="number"
                    value={character.Health}
                    onChange={(e) => {
                        const val = e.target.value === '' ? 1 : parseInt(e.target.value);
                        handleChange('Health', val);
                    }}
                    min="1"
                    max="999"
                />
            </div>
            <div className="form-group">
                <label htmlFor={`${label}-attack`}>Attack</label>
                <input
                    id={`${label}-attack`}
                    type="number"
                    value={character.Attack}
                    onChange={(e) => {
                        const val = e.target.value === '' ? 1 : parseInt(e.target.value);
                        handleChange('Attack', val);
                    }}
                    min="1"
                    max="100"
                />
            </div>
            <div className="form-group">
                <label htmlFor={`${label}-defense`}>Defense</label>
                <input
                    id={`${label}-defense`}
                    type="number"
                    value={character.Defense}
                    onChange={(e) => {
                        const val = e.target.value === '' ? 1 : parseInt(e.target.value);
                        handleChange('Defense', val);
                    }}
                    min="1"
                    max="100"
                />
            </div>
            <div className="form-group">
                <label htmlFor={`${label}-speed`}>Speed</label>
                <input
                    id={`${label}-speed`}
                    type="number"
                    value={character.Speed}
                    onChange={(e) => {
                        const val = e.target.value === '' ? 1 : parseInt(e.target.value);
                        handleChange('Speed', val);
                    }}
                    min="1"
                    max="100"
                />
            </div>
        </div>
    );
}
