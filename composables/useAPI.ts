// API Composable for backend communication
export const useAPI = () => {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase || 'http://localhost:4000/api/v1'

  // Helper function to handle API errors
  const handleError = (error: any) => {
    console.error('API Error:', error)
    if (error.response) {
      throw new Error(error.response._data?.error || 'API request failed')
    }
    throw error
  }

  return {
    // ========== HOUSE ENDPOINTS ==========
    
    /**
     * Get all houses
     */
    async getHouses(): Promise<House[]> {
      try {
        return await $fetch<House[]>(`${baseURL}/houses`)
      } catch (error) {
        handleError(error)
        return []
      }
    },

    /**
     * Get house by ID
     */
    async getHouseById(id: number): Promise<House | null> {
      try {
        return await $fetch<House>(`${baseURL}/houses/id/${id}`)
      } catch (error) {
        handleError(error)
        return null
      }
    },

    /**
     * Get houses by category
     */
    async getHousesByCategory(category: string): Promise<House[]> {
      try {
        return await $fetch<House[]>(`${baseURL}/houses/category/${category}`)
      } catch (error) {
        handleError(error)
        return []
      }
    },

    /**
     * Get latest houses
     */
    async getLatestHouses(limit: number = 10): Promise<House[]> {
      try {
        return await $fetch<House[]>(`${baseURL}/houses/latest`, {
          params: { limit }
        })
      } catch (error) {
        handleError(error)
        return []
      }
    },

    /**
     * Get houses by price range
     */
    async getHousesByPriceRange(minPrice: number, maxPrice: number): Promise<House[]> {
      try {
        return await $fetch<House[]>(`${baseURL}/houses/price`, {
          params: { minPrice, maxPrice }
        })
      } catch (error) {
        handleError(error)
        return []
      }
    },

    /**
     * Get house by address
     */
    async getHouseByAddress(address: string): Promise<House | null> {
      try {
        return await $fetch<House>(`${baseURL}/houses/address/${encodeURIComponent(address)}`)
      } catch (error) {
        handleError(error)
        return null
      }
    },

    // ========== AUTHENTICATION ENDPOINTS ==========

    /**
     * User login
     */
    async login(email: string, password: string): Promise<{ token: string }> {
      try {
        return await $fetch<{ token: string }>(`${baseURL}/login`, {
          method: 'POST',
          body: { email, password }
        })
      } catch (error) {
        handleError(error)
        throw error
      }
    },

    /**
     * User registration
     */
    async register(payload: RegisterPayload): Promise<{ message: string }> {
      try {
        return await $fetch<{ message: string }>(`${baseURL}/register`, {
          method: 'POST',
          body: payload
        })
      } catch (error) {
        handleError(error)
        throw error
      }
    },

    /**
     * Agent registration
     */
    async registerAgent(payload: RegisterPayload): Promise<{ message: string }> {
      try {
        return await $fetch<{ message: string }>(`${baseURL}/register/agent`, {
          method: 'POST',
          body: payload
        })
      } catch (error) {
        handleError(error)
        throw error
      }
    },

    // ========== REFERRAL ENDPOINTS ==========

    /**
     * Generate referral code (requires authentication)
     */
    async generateReferralCode(token: string): Promise<ReferralCode> {
      try {
        return await $fetch<ReferralCode>(`${baseURL}/referral/generate`, {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${token}`
          }
        })
      } catch (error) {
        handleError(error)
        throw error
      }
    },

    /**
     * Get referral statistics (requires authentication)
     */
    async getReferralStats(token: string): Promise<ReferralStats> {
      try {
        return await $fetch<ReferralStats>(`${baseURL}/referral/stats`, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        })
      } catch (error) {
        handleError(error)
        throw error
      }
    }
  }
}
