// Middleware to protect agent-only routes
export default defineNuxtRouteMiddleware((to, from) => {
    const { isAuthenticated, isAgent, isAdmin } = useAuth()

    // If not authenticated, redirect to login
    if (!isAuthenticated.value) {
        return navigateTo('/login')
    }

    // If authenticated but not an agent or admin, redirect to home
    if (!isAgent.value && !isAdmin.value) {
        return navigateTo('/')
    }
})
