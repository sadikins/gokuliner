<script setup>


  import NavbarComponent from './components/NavbarComponent.vue';



</script>

<template>
  <div class="flex h-screen">
    <!-- Navbar Aside -->
    <NavbarComponent />

    <!-- Main Contents -->
    <main class="flex-1 p-8 overflow-y-auto">
      <router-view />
    </main>

  </div>
</template>


