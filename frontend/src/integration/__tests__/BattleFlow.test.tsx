/// <reference types="jest" />
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import App from '../../App'

const mockFetch = jest.fn()
window.fetch = mockFetch

describe('Battle Flow Integration', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('handles API errors gracefully', async () => {
    // Mock failed API response
    mockFetch.mockImplementationOnce(() =>
      Promise.resolve({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error'
      })
    )

    render(<App />)

    const char1NameInput = screen.getByLabelText('Name', { selector: 'input[id="Character 1-name"]' })
    fireEvent.change(char1NameInput, { target: { value: 'Test Warrior' } })

    const createButton = screen.getByText('Create Battle')
    fireEvent.click(createButton)

    // Wait for error message
    await waitFor(() => {
      expect(screen.getByText('response.text is not a function')).toBeInTheDocument()
    })
  })
})
