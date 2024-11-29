import '@testing-library/jest-dom'
import 'whatwg-fetch'

// Mock fetch globally
global.fetch = jest.fn(() =>
  Promise.resolve(new Response(JSON.stringify({}), {
    status: 200,
    headers: { 'Content-type': 'application/json' }
  }))
) as jest.Mock

beforeEach(() => {
  (global.fetch as jest.Mock).mockClear()
})
