import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { BattleView } from '../BattleView'
import { mockBattle, mockCharacter1 } from '../../test/testUtils'
import { Battle } from '../../types'

describe('BattleView', () => {
  it('renders character names', () => {
    render(
      <BattleView
        battle={mockBattle}
        onAction={jest.fn()}
      />
    )

    expect(screen.getByText('Test Warrior')).toBeInTheDocument()
    expect(screen.getByText('Test Mage')).toBeInTheDocument()
  })

  it('handles actions', async () => {
    const onAction = jest.fn(() => Promise.resolve())
    render(
      <BattleView
        battle={mockBattle}
        onAction={onAction}
      />
    )

    // Find and click the first action button for the active character
    const buttons = screen.getAllByTestId('basic-attack')
    const activeCharacterButton = buttons.find(button => !button.hasAttribute('disabled'))
    expect(activeCharacterButton).toBeTruthy()
    fireEvent.click(activeCharacterButton!)

    await waitFor(() => {
      expect(onAction).toHaveBeenCalled()
    })
  })

  it('shows winner when battle is complete', () => {
    const completedBattle: Battle = {
      ...mockBattle,
      State: 'COMPLETE',
      Winner: mockCharacter1
    }
    render(<BattleView battle={completedBattle} onAction={jest.fn()} />)
    
    expect(screen.getByText(`Winner: ${mockCharacter1.Name}`)).toBeInTheDocument()
  })
})
