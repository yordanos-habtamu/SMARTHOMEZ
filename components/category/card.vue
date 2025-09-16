<template>
  <nuxt-link :to="`/categories/${category.name}`" class="relative group overflow-hidden rounded-lg block">
    <!-- Background Image -->
    <img 
      :src="category.image" 
      alt="Category image" 
      class="w-full h-60 object-cover transition-transform duration-700 group-hover:scale-110"
    />

    <!-- Centered Text Overlay -->
    <div class="absolute inset-0 bg-black/40 group-hover:bg-black/60 transition-all duration-700 flex items-center justify-center">
      <h3 class="text-xl group-hover:text-2xl transition-all duration-700 font-bold text-white text-center">
        {{ category.name }}
      </h3>
    </div>
  </nuxt-link>
</template>

<script setup lang="ts">
defineProps<{
  category: Category
}>()
</script>
