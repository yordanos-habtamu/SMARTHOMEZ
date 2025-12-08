<template>
    <section class="container py-12">
        <h2 class="text-3xl text-center font-bold text-gray-800 mb-6">
            Recently Listed Properties
        </h2>

        <div v-if="pending" class="flex justify-center py-12">
            <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-amber-500"></div>
        </div>

        <div v-else-if="error" class="text-center text-red-600 py-12">
            Failed to load properties. Please try again later.
        </div>

        <div v-else-if="properties && properties.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <property-card v-for="property in properties" :key="property.id" :property="property" />
        </div>

        <div v-else class="text-center text-gray-600 py-12">
            No recently listed properties.
        </div>
    </section>
</template>

<script setup lang="ts">
const api = useAPI()

// Fetch latest houses from API (limit to 8 for recent section)
const { data: properties, pending, error } = await useAsyncData(
    'recently-listed',
    () => api.getLatestHouses(8)
)
</script>