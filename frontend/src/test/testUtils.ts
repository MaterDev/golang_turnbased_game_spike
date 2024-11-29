import { Character, Battle } from '../types'

export const mockCharacter1: Character = {
  ID: 'char1-id',
  Name: 'Test Warrior',
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
}

export const mockCharacter2: Character = {
  ID: 'char2-id',
  Name: 'Test Mage',
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
}

export const mockBattle: Battle = {
  ID: 'battle-id',
  Character1: mockCharacter1,
  Character2: mockCharacter2,
  State: 'ACTIVE',
  Round: 1,
  Winner: undefined,
}
