import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from './auth'
import * as authService from '@/services/auth'
import * as tokenService from '@/services/token'

// Mock services
vi.mock('@/services/auth')
vi.mock('@/services/token')

describe('Auth Store', () => {
  let authStore

  beforeEach(() => {
    setActivePinia(createPinia())
    authStore = useAuthStore()

    // Clear all mocks
    vi.clearAllMocks()

    // Default mock implementations
    vi.mocked(tokenService.isTokenValid).mockReturnValue(true)
    vi.mocked(tokenService.getUserFromToken).mockReturnValue(null)
  })

  describe('Initial State', () => {
    it('should have correct initial state', () => {
      expect(authStore.user).toBeNull()
      expect(authStore.loading).toBe(false)
      expect(authStore.error).toBeNull()
      expect(authStore.isAuthenticated).toBe(false)
    })
  })

  describe('Login', () => {
    it('should login successfully with valid credentials', async () => {
      const mockResponse = {
        token: 'mock-jwt-token',
        expires_in: 3600,
        admin: {
          id: 1,
          email: 'admin@example.com',
          role: 'system_admin'
        }
      }

      vi.mocked(authService.login).mockResolvedValue(mockResponse)

      const result = await authStore.login('admin@example.com', 'password123')

      expect(result).toBe(true)
      expect(authStore.user).toEqual({
        adminId: 1,
        email: 'admin@example.com',
        role: 'system_admin'
      })
      expect(authStore.error).toBeNull()
      expect(tokenService.saveToken).toHaveBeenCalledWith('mock-jwt-token', 3600)
    })

    it('should handle login failure', async () => {
      const errorMessage = 'Invalid email or password'
      vi.mocked(authService.login).mockRejectedValue(new Error(errorMessage))

      const result = await authStore.login('admin@example.com', 'wrongpassword')

      expect(result).toBe(false)
      expect(authStore.user).toBeNull()
      expect(authStore.error).toBe(errorMessage)
    })

    it('should set loading state during login', async () => {
      let resolveFn
      const promise = new Promise(resolve => {
        resolveFn = resolve
      })

      vi.mocked(authService.login).mockReturnValue(promise)

      const loginPromise = authStore.login('admin@example.com', 'password123')

      expect(authStore.loading).toBe(true)

      resolveFn({
        token: 'token',
        expires_in: 3600,
        admin: { id: 1, email: 'admin@example.com', role: 'system_admin' }
      })

      await loginPromise

      expect(authStore.loading).toBe(false)
    })
  })

  describe('Logout', () => {
    beforeEach(async () => {
      // Login first
      vi.mocked(authService.login).mockResolvedValue({
        token: 'token',
        expires_in: 3600,
        admin: { id: 1, email: 'admin@example.com', role: 'system_admin' }
      })
      await authStore.login('admin@example.com', 'password123')
    })

    it('should logout successfully', async () => {
      vi.mocked(authService.logout).mockResolvedValue()

      await authStore.logout()

      expect(authStore.user).toBeNull()
      expect(authStore.error).toBeNull()
      expect(tokenService.removeToken).toHaveBeenCalled()
    })

    it('should clear user state even if logout API fails', async () => {
      vi.mocked(authService.logout).mockRejectedValue(new Error('API error'))

      await authStore.logout()

      expect(authStore.user).toBeNull()
      expect(tokenService.removeToken).toHaveBeenCalled()
    })
  })

  describe('Restore User', () => {
    it('should restore user from valid token', async () => {
      const tokenUser = {
        adminId: 1,
        email: 'admin@example.com',
        role: 'system_admin'
      }

      const apiResponse = {
        admin: {
          id: 1,
          email: 'admin@example.com',
          role: 'system_admin'
        }
      }

      vi.mocked(tokenService.isTokenValid).mockReturnValue(true)
      vi.mocked(tokenService.getUserFromToken).mockReturnValue(tokenUser)
      vi.mocked(authService.getCurrentUser).mockResolvedValue(apiResponse)

      const result = await authStore.restoreUser()

      expect(result).toBe(true)
      expect(authStore.user).toEqual(tokenUser)
    })

    it('should not restore user from invalid token', async () => {
      vi.mocked(tokenService.isTokenValid).mockReturnValue(false)

      const result = await authStore.restoreUser()

      expect(result).toBe(false)
      expect(authStore.user).toBeNull()
    })

    it('should clear token if API validation fails', async () => {
      vi.mocked(tokenService.isTokenValid).mockReturnValue(true)
      vi.mocked(authService.getCurrentUser).mockRejectedValue(new Error('Unauthorized'))

      const result = await authStore.restoreUser()

      expect(result).toBe(false)
      expect(authStore.user).toBeNull()
      expect(tokenService.removeToken).toHaveBeenCalled()
    })
  })

  describe('Computed Properties', () => {
    it('should correctly compute isAuthenticated', async () => {
      expect(authStore.isAuthenticated).toBe(false)

      // Login
      vi.mocked(authService.login).mockResolvedValue({
        token: 'token',
        expires_in: 3600,
        admin: { id: 1, email: 'admin@example.com', role: 'system_admin' }
      })
      await authStore.login('admin@example.com', 'password123')

      expect(authStore.isAuthenticated).toBe(true)
    })

    it('should correctly compute isSystemAdmin', async () => {
      vi.mocked(authService.login).mockResolvedValue({
        token: 'token',
        expires_in: 3600,
        admin: { id: 1, email: 'admin@example.com', role: 'system_admin' }
      })
      await authStore.login('admin@example.com', 'password123')

      expect(authStore.isSystemAdmin).toBe(true)
      expect(authStore.isAuctioneer).toBe(false)
    })

    it('should correctly compute isAuctioneer', async () => {
      vi.mocked(authService.login).mockResolvedValue({
        token: 'token',
        expires_in: 3600,
        admin: { id: 1, email: 'auctioneer@example.com', role: 'auctioneer' }
      })
      await authStore.login('auctioneer@example.com', 'password123')

      expect(authStore.isAuctioneer).toBe(true)
      expect(authStore.isSystemAdmin).toBe(false)
    })
  })

  describe('Clear Error', () => {
    it('should clear error', async () => {
      // Trigger an error
      vi.mocked(authService.login).mockRejectedValue(new Error('Login failed'))
      await authStore.login('admin@example.com', 'wrong')

      expect(authStore.error).toBe('Login failed')

      // Clear error
      authStore.clearError()

      expect(authStore.error).toBeNull()
    })
  })
})
