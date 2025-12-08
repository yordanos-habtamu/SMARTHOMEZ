<template>
  <div v-if="pending" class="flex justify-center items-center min-h-screen">
    <div class="text-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-amber-500 mx-auto mb-4"></div>
      <p class="text-gray-600">Loading property...</p>
    </div>
  </div>

  <div v-else-if="error" class="flex justify-center items-center min-h-screen">
    <div class="text-center">
      <h2 class="text-2xl font-bold text-red-600 mb-4">Property Not Found</h2>
      <p class="text-gray-600 mb-6">{{ error.message || 'Unable to load property details' }}</p>
      <nuxt-link to="/listing" class="bg-amber-500 text-white px-6 py-2 rounded-lg hover:bg-amber-600">
        Back to Listings
      </nuxt-link>
    </div>
  </div>

  <div v-else-if="property">
    <!-- Hero Section -->
    <section class="relative h-[60vh] flex items-center justify-center rounded-lg">
      <img
        :src="property.imgUrl"
        :alt="property.address"
        class="absolute inset-0 w-full h-full object-cover"
      />
      <div class="absolute inset-0 bg-black opacity-50" />
    </section>

    <div class="relative container flex flex-col items-center justify-center w-full z-20">
      <section class="bg-white rounded-lg shadow-lg -mt-20 p-4 lg:p-8 mb-12 max-w-5xl">
        <h2 class="text-3xl font-bold text-center mb-2">{{ property.address }}</h2>

        <p class="text-gray-700 text-center text-lg">
          {{ property.description }}
        </p>

        <div class="text-gray-700 text-center mb-8">
          Listed {{ formatDate(property.createdAt) }}
        </div>

        <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
          <div class="bg-zinc-100 rounded-lg py-4">
            <div class="text-2xl font-bold text-center">
              {{ property.numBedrooms }}
            </div>
            <p class="text-center text-gray-700">Bedrooms</p>
          </div>
          
          <div class="bg-zinc-100 rounded-lg py-4">
            <div class="text-2xl font-bold text-center">
              {{ property.numBathrooms }}
            </div>
            <p class="text-center text-gray-700">Bathrooms</p>
          </div>

          <div class="bg-zinc-100 rounded-lg py-4">
            <div class="text-2xl font-bold text-center">
              {{ property.areaSqFt }} sq ft
            </div>
            <p class="text-center text-gray-700">Square Feet</p>
          </div>
          
          <div class="bg-zinc-100 rounded-lg py-4">
            <div class="text-2xl font-black text-center text-primary">
              ${{ property.price.toLocaleString() }}
            </div>
            <p class="text-center text-gray-700">Price</p>
          </div>
        </div>

        <!-- Additional Details -->
        <div class="mt-6 border-t pt-6">
          <div class="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span class="font-semibold">Category:</span>
              <span class="ml-2 text-gray-600">{{ property.category }}</span>
            </div>
            <div>
              <span class="font-semibold">Status:</span>
              <span class="ml-2" :class="property.isSold ? 'text-red-600' : 'text-green-600'">
                {{ property.isSold ? 'Sold' : 'Available' }}
              </span>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const api = useAPI()

// Validate route parameter
const propertyId = computed(() => {
  const id = Number(route.params.id)
  return isNaN(id) ? null : id
})

// Fetch property data from API
const { data: property, pending, error } = await useAsyncData(
  'property-' + route.params.id,
  async () => {
    if (!propertyId.value) {
      throw new Error('Invalid property ID')
    }
    const result = await api.getHouseById(propertyId.value)
    if (!result) {
      throw new Error('Property not found')
    }
    return result
  }
)

// Format date helper
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric' 
  })
}

// Set page meta
definePageMeta({
  validate: async (route) => {
    const id = Number(route.params.id)
    return !isNaN(id) && id > 0
  }
})
</script>