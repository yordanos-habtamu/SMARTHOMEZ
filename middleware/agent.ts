// Middleware to protect agent-only routes
export default defineNuxtRouteMiddleware((to, from) => {
    const { isAuthenticated, isAgent } = useAuth()

    // If not authenticated, redirect to login
    if (!isAuthenticated.value) {
        return navigateTo('/login')
    }

    // If authenticated but not an agent, redirect to home
    if (!isAgent.value) {
        return navigateTo('/')
    }
})
