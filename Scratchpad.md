# Considerations (To Look Into)

- Directory structure for testing:
  - Do i put all the tests in one place, or try to co-locate them with what they belong to.
  - For character `struct`
      - Map out the relationship diagram between the entities I am using in this project.
      - Get a clearer idea of the shape and relationships between the data.
      - []Abilities
          - Do I want to have a fixed number of abilities
          - Impose hard cap on the number of abilities

- Game Machine
  - When game is initiated attacks and status effects will manipulate a copy of the character data structure, which is stored in memory along side its original.
  - Character data structure and state management.
    - Make each round of game into a new index in a slice of state objects- which includes snapshots of resultant character state at the end of game rounds.
    - A linked list of character states. (Circular doubly linked list with a reference to Head in each node.)

3 Possible Members of the Character Data Structure family:
  - The original Character data struct (Immutable)
  - A mutable copy, to which damage and status effect apply
  - A slice of snapshots from mutable copy - per end-of-round
