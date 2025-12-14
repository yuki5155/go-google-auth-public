<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<template>
  <header class="header">
    <div class="header-content">
      <h1 class="header-title">Go Google Auth</h1>
      <nav class="header-nav">
        <RouterLink to="/">Session Test</RouterLink>
        <RouterLink to="/about">About</RouterLink>
      </nav>
    </div>
  </header>
</template>
