/// <reference types="jest" />
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import App from '../../App'
import { mockBattle } from '../../test/testUtils'

const mockFetch = jest.fn()
window.fetch = mockFetch

describe('App', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('shows initial battle creation form', () => {
    render(<App />)
    expect(screen.getByText('Create New Battle')).toBeInTheDocument()
  })

  it('creates a new battle', async () => {
    mockFetch.mockImplementationOnce(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve(mockBattle)
      })
    )

    render(<App />)

    // Fill in character forms using more specific selectors
    const char1NameInput = screen.getByLabelText('Name', { selector: 'input[id="Character 1-name"]' })
    fireEvent.change(char1NameInput, { target: { value: 'Test Warrior' } })

    // Click create battle
    const createButton = screen.getByText('Create Battle')
    fireEvent.click(createButton)

    // Wait for battle to be created
    await waitFor(() => {
      expect(mockFetch).toHaveBeenCalledTimes(1)
    })
  })

  it('handles API errors gracefully', async () => {
    mockFetch.mockImplementationOnce(() =>
      Promise.reject(new Error('API Error'))
    )

    render(<App />)

    // Click create battle
    const createButton = screen.getByText('Create Battle')
    fireEvent.click(createButton)

    // Wait for error message - look for the actual error text
    await waitFor(() => {
      expect(screen.getByText('API Error')).toBeInTheDocument()
    })
  })

  it('resets state when clicking Start New Battle', async () => {
    // First create a battle with an error to show the reset button
    mockFetch.mockImplementationOnce(() =>
      Promise.reject(new Error('API Error'))
    )

    render(<App />)

    // Click create battle
    const createButton = screen.getByText('Create Battle')
    fireEvent.click(createButton)

    // Wait for the error state and reset button
    await waitFor(() => {
      expect(screen.getByText('API Error')).toBeInTheDocument()
    })

    // Now the reset button should be visible
    const resetButton = screen.getByRole('button', { name: 'Start New Battle' })
    fireEvent.click(resetButton)

    // Check if we're back to initial state
    expect(screen.getByText('Create New Battle')).toBeInTheDocument()
  })
})
