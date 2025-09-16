<template>
    <header class="absolute top-0 left-0 w-full border-b-gray-800 border-b z-20 p-6"> 
    <div class="container flex justify-between items-center">
        <h1 class="text-center text-lg text-amber-500 font-extralight">SmartHomez</h1>
       <nav>
        <ul class="flex space-x-4">
     
            <li v-for="item in navLinks" :key="item.name" :href="item.link" class="text-white">
               <nuxt-link :to="item.link"> {{ item.name }}</nuxt-link>
            </li>
      
        </ul>
       </nav>
    </div>
    </header>
</template>
<script setup lang='ts'> 
import {ref} from 'vue';
const navLinks =ref([{
    name:"Home",
    link:"/"
},
{
    name:"About Us",
    link:"/about"
},
{
    name:"Listing",
    link:"/listing"
},
{
    name:"Contact",
    link:"/contact"
}]);

</script>