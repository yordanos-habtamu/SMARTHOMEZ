// Backend-aligned type definitions for SMARTHOMEZ API

export { }
declare global {
    // House/Property type matching backend
    type House = {
        id: number                 // Backend uses uint
        address: string            // Changed from 'title'
        price: number              // Backend uses float64
        numBedrooms: number        // Aligned with backend
        numBathrooms: number       // Aligned with backend
        areaSqFt: number          // Changed from squareFeet
        category: string
        description: string
        imgUrl: string            // Changed from 'image', backend uses ImgURL
        isSold: boolean
        agentId: number
        createdAt: string         // ISO date string
        updatedAt: string         // ISO date string
    }

    // Category type (frontend only, no backend equivalent yet)
    type Category = {
        id: string
        name: string
        description: string
        image: string
    }

    // API Response types
    type ApiResponse<T> = {
        data?: T
        error?: string
        message?: string
    }

    // User type for authentication
    type User = {
        id: number
        firstName: string
        lastName: string
        email: string
        role: 'admin' | 'agent' | 'customer'
        createdAt: string
    }

    // Login/Register payload types
    type LoginPayload = {
        email: string
        password: string
    }

    type RegisterPayload = {
        firstName: string
        lastName: string
        email: string
        password: string
        contact: string
        DoB: string           // Format: YYYY-MM-DD
        sex: string
        role?: string
        referralCode?: string
    }

    // Referral types
    type ReferralCode = {
        id: number
        userId: number
        code: string
        shortCode: string
        qrCodeUrl: string
        visitCount: number
        signupCount: number
        createdAt: string
    }

    type ReferralStats = {
        code: string
        shortCode: string
        qrCodeUrl: string
        visitCount: number
        signupCount: number
    }
}