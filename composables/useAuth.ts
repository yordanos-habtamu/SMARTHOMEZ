import { useAPI } from "./useAPI"

// Authentication composable for managing user auth state
export const useAuth = () => {
    const token = useCookie('auth_token', { maxAge: 60 * 60 * 24 * 7 })
    const user = useState<User | null>('user', () => null)

    // Check if user is authenticated
    const isAuthenticated = computed(() => !!token.value)

    // Check if user is an agent
    const isAgent = computed(() => user.value?.role === 'agent')

    // Check if user is admin
    const isAdmin = computed(() => user.value?.role === 'admin')

    // Login function
    const login = async (email: string, password: string) => {
        const api = useAPI()
        const { token: authToken } = await api.login(email, password)
        token.value = authToken

        // Decode JWT to get user info (basic decode, not validation)
        // In production, you'd validate this on the backend
        if (authToken) {
            const payload = JSON.parse(atob(authToken.split('.')[1]))
            user.value = {
                id: payload.userId,
                role: payload.role,
                firstName: '',
                lastName: '',
                email: email,
                createdAt: ''
            }
        }
    }

    // Logout function
    const logout = () => {
        token.value = null
        user.value = null
        navigateTo('/login')
    }

    // Get user from token on mount
    const initUser = () => {
        if (token.value) {
            try {
                const payload = JSON.parse(atob(token.value.split('.')[1]))
                user.value = {
                    id: payload.userId,
                    role: payload.role,
                    firstName: '',
                    lastName: '',
                    email: '',
                    createdAt: ''
                }
            } catch (e) {
                console.error('Failed to parse token:', e)
                logout()
            }
        }
    }

    return {
        token,
        user,
        isAuthenticated,
        isAgent,
        isAdmin,
        login,
        logout,
        initUser
    }
}
