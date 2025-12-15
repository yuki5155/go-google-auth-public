<script setup lang="ts">
import AppHeader from './components/AppHeader.vue'
</script>

<template>
  <div id="app">
    <AppHeader />
    <main class="main-content">
      <RouterView />
    </main>
  </div>
</template>

<style>
.main-content {
  flex: 1;
  width: 100%;
  padding: 2rem;
}
</style>
