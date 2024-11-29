export type Character = {
    ID: string;
    Name: string;
    Health: number;
    Attack: number;
    Defense: number;
    Speed: number;
    StatusEffects: StatusEffect[];
    Abilities: Ability[];
};

export type StatusEffect = {
    Type: string;
    Duration: number;
    Potency: number;
};

export type Ability = {
    Name: string;
    Damage: number;
    CooldownMax: number;
    StatusEffect?: StatusEffect;
};

export type Battle = {
    ID: string;
    Character1: Character;
    Character2: Character;
    State: "PENDING" | "ACTIVE" | "COMPLETE";
    Winner?: Character;
    Round: number;
};

export type BattleAction = {
    CharacterID: string;
    AbilityIndex: number;
    TargetID: string;
};
