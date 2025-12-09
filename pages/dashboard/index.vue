<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Dashboard Header -->
    <div class="bg-white shadow">
      <div class="container mx-auto px-4 py-6">
        <div class="flex justify-between items-center">
          <div>
            <h1 class="text-3xl font-bold text-gray-800">{{ isAdmin ? 'Admin Dashboard' : 'Agent Dashboard' }}</h1>
            <p class="text-gray-600 mt-1">{{ isAdmin ? 'Manage system overview' : 'Manage your properties and referrals' }}</p>
          </div>
          <button
            @click="handleLogout"
            class="bg-red-500 text-white px-4 py-2 rounded-lg hover:bg-red-600"
          >
            Logout
          </button>
        </div>
      </div>
    </div>

    <!-- Dashboard Navigation -->
    <div class="container mx-auto px-4 py-6">
      <div class="flex space-x-4 border-b mb-6">
        <nuxt-link
          to="/dashboard"
          class="px-4 py-2 font-semibold border-b-2"
          :class="route.path === '/dashboard' ? 'border-amber-500 text-amber-600' : 'border-transparent text-gray-600 hover:text-gray-800'"
        >
          Overview
        </nuxt-link>
        <nuxt-link
          to="/dashboard/referrals"
          class="px-4 py-2 font-semibold border-b-2"
          :class="route.path === '/dashboard/referrals' ? 'border-amber-500 text-amber-600' : 'border-transparent text-gray-600 hover:text-gray-800'"
        >
          Referrals
        </nuxt-link>
      </div>

      <!-- Stats Cards -->
      <div v-if="pending" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-amber-500"></div>
      </div>

      <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
        <p>{{ errorMessage }}</p>
        <p class="text-sm mt-2">You may need to generate a referral code first.</p>
        <nuxt-link to="/dashboard/referrals" class="text-red-800 underline font-semibold">
          Go to Referrals →
        </nuxt-link>
      </div>

      <div v-else>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <div class="bg-white rounded-lg shadow p-6">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-gray-600 text-sm">Total Referral Visits</p>
                <p class="text-3xl font-bold text-gray-800 mt-2">
                  {{ stats?.visitCount || 0 }}
                </p>
              </div>
              <div class="bg-blue-100 rounded-full p-3">
                <svg class="w-8 h-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                </svg>
              </div>
            </div>
          </div>

          <div class="bg-white rounded-lg shadow p-6">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-gray-600 text-sm">Total Signups</p>
                <p class="text-3xl font-bold text-gray-800 mt-2">
                  {{ stats?.signupCount || 0 }}
                </p>
              </div>
              <div class="bg-green-100 rounded-full p-3">
                <svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"/>
                </svg>
              </div>
            </div>
          </div>

          <div class="bg-white rounded-lg shadow p-6">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-gray-600 text-sm">Conversion Rate</p>
                <p class="text-3xl font-bold text-gray-800 mt-2">
                  {{ conversionRate }}%
                </p>
              </div>
              <div class="bg-amber-100 rounded-full p-3">
                <svg class="w-8 h-8 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
                </svg>
              </div>
            </div>
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="bg-white rounded-lg shadow p-6">
          <h2 class="text-xl font-bold text-gray-800 mb-4">Quick Actions</h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <nuxt-link
              to="/dashboard/referrals"
              class="flex items-center p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition"
            >
              <div class="bg-amber-100 rounded-full p-3 mr-4">
                <svg class="w-6 h-6 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z"/>
                </svg>
              </div>
              <div>
                <h3 class="font-semibold text-gray-800">View Referrals</h3>
                <p class="text-sm text-gray-600">Check your referral stats and QR code</p>
              </div>
            </nuxt-link>

            <nuxt-link
              to="/"
              class="flex items-center p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition"
            >
              <div class="bg-blue-100 rounded-full p-3 mr-4">
                <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
                </svg>
              </div>
              <div>
                <h3 class="font-semibold text-gray-800">Browse Properties</h3>
                <p class="text-sm text-gray-600">View all available properties</p>
              </div>
            </nuxt-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAPI } from '~/composables/useAPI'

const route = useRoute()
const { logout, token, isAdmin } = useAuth()
const api = useAPI()

// Fetch referral stats
const { data: stats, pending, error } = await useAsyncData(
  'referral-stats-overview',
  async () => {
    if (!token.value) throw new Error('Not authenticated')
    try {
      return await api.getReferralStats(token.value)
    } catch (e) {
      // User may not have generated a referral code yet
      return null
    }
  }
)

const errorMessage = computed(() => {
  if (error.value) return 'Failed to load referral stats'
  return ''
})

// Calculate conversion rate
const conversionRate = computed(() => {
  if (!stats.value || stats.value.visitCount === 0) return 0
  return ((stats.value.signupCount / stats.value.visitCount) * 100).toFixed(1)
})

const handleLogout = () => {
  logout()
}

// Protect this page with agent middleware
definePageMeta({
  middleware: 'agent'
})
</script>
