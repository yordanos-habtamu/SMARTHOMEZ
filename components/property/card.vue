<template>
 <nuxt-link :to="`/property/${property.id}`" class="max-w-sm mx-auto bg-white shadow-lg rounded-lg overflow-hidden hover:shadow-xl transition-shadow">
   <img 
     :src="property.imgUrl" 
     :alt="property.address"
     class="w-full h-48 object-cover"
   >
   <div class="p-6">
    <h2 class="text-lg font-bold text-gray-800 mb-2">{{ property.address }}</h2>
    <p class="text-gray-600 text-sm mb-4 line-clamp-2">{{ property.description }}</p>
    
    <div class="flex justify-between items-center mb-3">
        <span class="font-bold text-amber-500 text-lg">
            ${{ property.price.toLocaleString() }}
        </span>
        <span 
          :class="property.isSold ? 'bg-red-100 text-red-600' : 'bg-green-100 text-green-600'"
          class="text-xs px-2 py-1 rounded"
        >
            {{ property.isSold ? 'Sold' : 'Available' }}
        </span>
    </div>
    
    <div class="flex justify-between items-center text-gray-600 text-sm">
        <span>{{ property.numBedrooms }} BR</span>
        <span>|</span>
        <span>{{ property.numBathrooms }} BA</span>
        <span>|</span>
        <span>{{ property.areaSqFt }} sqft</span>
    </div>
   </div>
 </nuxt-link>
</template>

<script setup lang="ts">
defineProps<{
    property: House
}>()
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>