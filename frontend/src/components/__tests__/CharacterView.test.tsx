import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { CharacterView } from '../BattleView'
import { mockCharacter1, mockCharacter2 } from '../../test/testUtils'
import { Character } from '../../types'

const defaultProps = {
  character: mockCharacter1,
  opponent: mockCharacter2,
  onAction: jest.fn().mockResolvedValue(undefined),
  disabled: false,
  isCurrentTurn: true,
}

describe('CharacterView', () => {
  beforeEach(() => {
    defaultProps.onAction.mockClear()
  })

  it('renders character information correctly', () => {
    render(<CharacterView {...defaultProps} />)
    
    expect(screen.getByText(mockCharacter1.Name)).toBeInTheDocument()
    expect(screen.getByText(`Health: ${mockCharacter1.Health}`)).toBeInTheDocument()
    expect(screen.getByText(`Attack: ${mockCharacter1.Attack}`)).toBeInTheDocument()
    expect(screen.getByText(`Defense: ${mockCharacter1.Defense}`)).toBeInTheDocument()
    expect(screen.getByText(`Speed: ${mockCharacter1.Speed}`)).toBeInTheDocument()
  })

  it('disables actions when not current turn', () => {
    render(<CharacterView {...defaultProps} isCurrentTurn={false} />)
    
    const actionButtons = screen.getAllByRole('button')
    actionButtons.forEach(button => {
      expect(button).toBeDisabled()
    })
  })

  it('handles target selection and basic attack', async () => {
    render(<CharacterView {...defaultProps} />)
    
    const select = screen.getByRole('combobox')
    fireEvent.change(select, { target: { value: 'self' } })
    
    const basicAttackButton = screen.getByTestId('basic-attack')
    fireEvent.click(basicAttackButton)
    
    await waitFor(() => {
      expect(defaultProps.onAction).toHaveBeenCalledWith({
        CharacterID: mockCharacter1.ID,
        AbilityIndex: 0,
        TargetID: mockCharacter1.ID,
      })
    })
  })

  it('handles special attack', async () => {
    render(<CharacterView {...defaultProps} />)
    
    const specialAttackButton = screen.getByTestId('special-attack')
    fireEvent.click(specialAttackButton)
    
    await waitFor(() => {
      expect(defaultProps.onAction).toHaveBeenCalledWith({
        CharacterID: mockCharacter1.ID,
        AbilityIndex: 1,
        TargetID: mockCharacter2.ID,
      })
    })
  })

  it('shows loading state during action submission', async () => {
    defaultProps.onAction.mockImplementation(() => new Promise(resolve => setTimeout(resolve, 100)))
    render(<CharacterView {...defaultProps} />)
    
    const basicAttackButton = screen.getByTestId('basic-attack')
    fireEvent.click(basicAttackButton)
    
    expect(basicAttackButton).toBeDisabled()
    
    await waitFor(() => {
      expect(basicAttackButton).not.toBeDisabled()
    }, { timeout: 150 })
  })

  it('handles status effects display', () => {
    const characterWithStatus: Character = {
      ...mockCharacter1,
      StatusEffects: [
        { Type: 'BURNING', Duration: 2, Potency: 5 },
      ],
    }
    
    render(<CharacterView {...defaultProps} character={characterWithStatus} />)
    expect(screen.getByText(/Status Effects: BURNING/)).toBeInTheDocument()
  })

  it('disables actions when character is disabled', () => {
    render(<CharacterView {...defaultProps} disabled={true} />)
    
    const actionButtons = screen.getAllByRole('button')
    actionButtons.forEach(button => {
      expect(button).toBeDisabled()
    })
  })

  it('shows current turn indicator', () => {
    const { container } = render(<CharacterView {...defaultProps} isCurrentTurn={true} />)
    expect(container.querySelector('.current-turn')).toBeInTheDocument()
  })

  it('handles missing character data', () => {
    const props = {
      ...defaultProps,
      character: undefined as any,
    }
    
    render(<CharacterView {...props} />)
    expect(screen.getByText('Loading character data...')).toBeInTheDocument()
  })
})
