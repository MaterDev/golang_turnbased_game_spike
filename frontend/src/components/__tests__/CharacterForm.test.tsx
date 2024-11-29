import { render, screen, fireEvent } from '@testing-library/react'
import { CharacterForm } from '../CharacterForm'
import { Character } from '../../types'

const mockCharacter: Character = {
  ID: '',
  Name: 'Test Character',
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
    }
  ],
}

describe('CharacterForm', () => {
  const mockOnSubmit = jest.fn()

  beforeEach(() => {
    mockOnSubmit.mockClear()
  })

  it('renders with initial values', () => {
    render(
      <CharacterForm
        label="Test Form"
        defaultCharacter={mockCharacter}
        onChange={() => {}}
      />
    )

    expect(screen.getByLabelText('Name')).toHaveValue('Test Character')
    expect(Number((screen.getByLabelText('Health') as HTMLInputElement).value)).toBe(100)
    expect(Number((screen.getByLabelText('Attack') as HTMLInputElement).value)).toBe(15)
    expect(Number((screen.getByLabelText('Defense') as HTMLInputElement).value)).toBe(10)
    expect(Number((screen.getByLabelText('Speed') as HTMLInputElement).value)).toBe(8)
  })

  it('calls onChange when values change', () => {
    const onChange = jest.fn()
    render(
      <CharacterForm
        label="Test Form"
        defaultCharacter={mockCharacter}
        onChange={onChange}
      />
    )

    const nameInput = screen.getByLabelText('Name')
    fireEvent.change(nameInput, { target: { value: 'New Name' } })
    expect(onChange).toHaveBeenCalledWith(expect.objectContaining({
      Name: 'New Name'
    }))

    const healthInput = screen.getByLabelText('Health') as HTMLInputElement
    fireEvent.change(healthInput, { target: { value: '150' } })
    expect(onChange).toHaveBeenCalledWith(expect.objectContaining({
      Health: 150
    }))
  })

  it('clamps numeric values within bounds', () => {
    render(
      <CharacterForm
        label="Test Form"
        defaultCharacter={mockCharacter}
        onChange={() => {}}
      />
    )

    // Health (1-999)
    const healthInput = screen.getByLabelText('Health') as HTMLInputElement
    fireEvent.change(healthInput, { target: { value: '1000' } })
    expect(Number(healthInput.value)).toBe(999)
    fireEvent.change(healthInput, { target: { value: '0' } })
    expect(Number(healthInput.value)).toBe(1)

    // Attack (1-100)
    const attackInput = screen.getByLabelText('Attack') as HTMLInputElement
    fireEvent.change(attackInput, { target: { value: '101' } })
    expect(Number(attackInput.value)).toBe(100)
    fireEvent.change(attackInput, { target: { value: '0' } })
    expect(Number(attackInput.value)).toBe(1)
  })

  it('handles empty input values', () => {
    render(
      <CharacterForm
        label="Test Form"
        defaultCharacter={mockCharacter}
        onChange={() => {}}
      />
    )

    const healthInput = screen.getByLabelText('Health') as HTMLInputElement
    fireEvent.change(healthInput, { target: { value: '' } })
    expect(Number(healthInput.value)).toBe(1)
  })
})
