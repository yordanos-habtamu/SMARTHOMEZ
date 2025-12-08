<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Dashboard Header -->
    <div class="bg-white shadow">
      <div class="container mx-auto px-4 py-6">
        <div class="flex justify-between items-center">
          <div>
            <h1 class="text-3xl font-bold text-gray-800">Referral System</h1>
            <p class="text-gray-600 mt-1">Share your referral link and earn rewards</p>
          </div>
          <nuxt-link
            to="/dashboard"
            class="text-amber-600 hover:text-amber-700 font-semibold"
          >
            ← Back to Dashboard
          </nuxt-link>
        </div>
      </div>
    </div>

    <div class="container mx-auto px-4 py-8">
      <div v-if="pending" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-amber-500"></div>
      </div>

      <div v-else-if="error" class="text-center py-12">
        <p class="text-red-600 mb-4">{{ error.message || 'Failed to load referral data' }}</p>
        <button
          @click="generateReferral"
          class="bg-amber-500 text-white px-6 py-2 rounded-lg hover:bg-amber-600"
        >
          Generate Referral Code
        </button>
      </div>

      <div v-else-if="!referralData && !generating">
        <div class="bg-white rounded-lg shadow p-8 text-center">
          <h2 class="text-2xl font-bold text-gray-800 mb-4">No Referral Code Yet</h2>
          <p class="text-gray-600 mb-6">Generate your unique referral code to start earning rewards</p>
          <button
            @click="generateReferral"
            class="bg-amber-500 text-white px-6 py-3 rounded-lg hover:bg-amber-600 font-semibold"
          >
            Generate My Referral Code
          </button>
        </div>
      </div>

      <div v-else-if="generating" class="text-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-amber-500 mx-auto mb-4"></div>
        <p class="text-gray-600">Generating your referral code...</p>
      </div>

      <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- QR Code Card -->
        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-lg font-bold text-gray-800 mb-4">Your QR Code</h3>
          <div class="bg-gray-100 rounded-lg p-4 flex items-center justify-center">
            <img
              v-if="referralData.qrCodeUrl"
              :src="referralData.qrCodeUrl"
              alt="Referral QR Code"
              class="w-48 h-48"
            />
            <div v-else class="text-gray-500">QR Code not available</div>
          </div>
          <p class="text-sm text-gray-600 text-center mt-4">
            Share this QR code with potential clients
          </p>
        </div>

        <!-- Referral Links Card -->
        <div class="bg-white rounded-lg shadow p-6 lg:col-span-2">
          <h3 class="text-lg font-bold text-gray-800 mb-4">Referral Links</h3>

          <!-- Full Referral Code -->
          <div class="mb-4">
            <label class="block text-sm font-semibold text-gray-700 mb-2">
              Referral Code
            </label>
            <div class="flex">
              <input
                :value="referralData.code"
                readonly
                class="flex-1 px-4 py-2 border border-gray-300 rounded-l-lg bg-gray-50"
              />
              <button
                @click="copyToClipboard(referralData.code, 'code')"
                class="bg-amber-500 text-white px-4 py-2 rounded-r-lg hover:bg-amber-600"
              >
                {{ copiedCode ? '✓ Copied' : 'Copy' }}
              </button>
            </div>
          </div>

          <!-- Short Link -->
          <div class="mb-4">
            <label class="block text-sm font-semibold text-gray-700 mb-2">
              Short Link
            </label>
            <div class="flex">
              <input
                :value="shortLink"
                readonly
                class="flex-1 px-4 py-2 border border-gray-300 rounded-l-lg bg-gray-50"
              />
              <button
                @click="copyToClipboard(shortLink, 'link')"
                class="bg-amber-500 text-white px-4 py-2 rounded-r-lg hover:bg-amber-600"
              >
                {{ copiedLink ? '✓ Copied' : 'Copy' }}
              </button>
            </div>
          </div>

          <!-- Share Buttons -->
          <div class="mt-6">
            <p class="text-sm font-semibold text-gray-700 mb-3">Share via:</p>
            <div class="flex space-x-3">
              <a
                :href="`https://wa.me/?text=Join SMARTHOMEZ using my referral link: ${shortLink}`"
                target="_blank"
                class="bg-green-500 text-white px-4 py-2 rounded-lg hover:bg-green-600 flex items-center"
              >
                <svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.890-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413z"/>
                </svg>
                WhatsApp
              </a>
              <a
                :href="`mailto:?subject=Join SMARTHOMEZ&body=Use my referral link: ${shortLink}`"
                class="bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600 flex items-center"
              >
                <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
                </svg>
                Email
              </a>
            </div>
          </div>
        </div>

        <!-- Stats Card -->
        <div class="bg-white rounded-lg shadow p-6 lg:col-span-3">
          <h3 class="text-lg font-bold text-gray-800 mb-4">Referral Statistics</h3>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="border border-gray-200 rounded-lg p-4 text-center">
              <p class="text-3xl font-bold text-blue-600">{{ referralData.visitCount }}</p>
              <p class="text-sm text-gray-600 mt-1">Total Visits</p>
            </div>
            <div class="border border-gray-200 rounded-lg p-4 text-center">
              <p class="text-3xl font-bold text-green-600">{{ referralData.signupCount }}</p>
              <p class="text-sm text-gray-600 mt-1">Total Signups</p>
            </div>
            <div class="border border-gray-200 rounded-lg p-4 text-center">
              <p class="text-3xl font-bold text-amber-600">{{ conversionRate }}%</p>
              <p class="text-sm text-gray-600 mt-1">Conversion Rate</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const config = useRuntimeConfig()
const { token } = useAuth()
const api = useAPI()

const generating = ref(false)
const copiedCode = ref(false)
const copiedLink = ref(false)

// Fetch referral data
const { data: referralData, pending, error, refresh } = await useAsyncData(
  'referral-data',
  async () => {
    if (!token.value) return null
    try {
      return await api.getReferralStats(token.value)
    } catch (e) {
      // If no referral exists yet, return null
      return null
    }
  }
)

// Generate short link
const shortLink = computed(() => {
  if (!referralData.value?.shortCode) return ''
  const baseUrl = config.public.apiBase.replace('/api/v1', '')
  return `${baseUrl}/r/${referralData.value.shortCode}`
})

// Calculate conversion rate
const conversionRate = computed(() => {
  if (!referralData.value || referralData.value.visitCount === 0) return 0
  return ((referralData.value.signupCount / referralData.value.visitCount) * 100).toFixed(1)
})

// Generate new referral code
const generateReferral = async () => {
  if (!token.value) return
  generating.value = true
  try {
    await api.generateReferralCode(token.value)
    await refresh()
  } catch (e: any) {
    alert('Failed to generate referral code: ' + (e.message || 'Unknown error'))
  } finally {
    generating.value = false
  }
}

// Copy to clipboard
const copyToClipboard = async (text: string, type: 'code' | 'link') => {
  try {
    await navigator.clipboard.writeText(text)
    if (type === 'code') {
      copiedCode.value = true
      setTimeout(() => copiedCode.value = false, 2000)
    } else {
      copiedLink.value = true
      setTimeout(() => copiedLink.value = false, 2000)
    }
  } catch (e) {
    alert('Failed to copy to clipboard')
  }
}

// Protect this page with agent middleware
definePageMeta({
  middleware: 'agent'
})
</script>
