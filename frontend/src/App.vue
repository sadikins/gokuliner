<!-- src/App.vue -->
<template>
  <div id="app">
    <template v-if="$route.meta.requiresAuth">

        <router-view />

    </template>
    <template v-else>
      <router-view />
    </template>
  </div>
</template>

<script>

export default {
  name: 'App',
  components: {
  },
};
</script>

<style>
/* Global styles for App.vue or import main.css here */
/* Pastikan main.css sudah diimport di main.js */
</style>