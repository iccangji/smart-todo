<template>
  <UCard>
    <template #header>
      <h3 class="text-md font-semibold text-gray-900 dark:text-white">
        {{ title }}
      </h3>
    </template>
    <div class="text-3xl font-bold text-primary dark:text-primary-400">
      {{ value }}
    </div>
  </UCard>
</template>

<script setup lang="ts">
defineProps({
  title: {
    type: String,
    required: true,
  },
  value: {
    type: [String, Number],
    required: true,
  },
});
</script>
