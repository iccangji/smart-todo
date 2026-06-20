<template>
  <AppLayout>
    <div class="space-y-6">
      <h1 class="text-2xl font-bold">Dashboard</h1>

      <!-- Recommendation -->
      <div class="p-4 border rounded-xl">
        <p class="text-gray-500 text-sm">Daily Recommendation</p>
        <p class="text-lg font-medium mt-2">
          {{ dashboard.recommendation?.message }}
        </p>
      </div>

      <!-- Summary Cards -->
      <div class="grid grid-cols-4 gap-4">
        <div class="card">
          <p>Total</p>
          <h2>{{ dashboard.summary?.total }}</h2>
        </div>

        <div class="card">
          <p>Completed</p>
          <h2>{{ dashboard.summary?.completed_count }}</h2>
        </div>

        <div class="card">
          <p>Pending</p>
          <h2>{{ dashboard.summary?.pending_count }}</h2>
        </div>

        <div class="card">
          <p>Completion Rate</p>
          <h2>{{ dashboard.summary?.completion_rate }}%</h2>
        </div>
      </div>

      <!-- Priority -->
      <div class="p-4 border rounded-xl">
        <h3 class="font-semibold mb-2">Priority Breakdown</h3>

        <pre
          >{{ dashboard.summary?.pending_priority_count }}
        </pre>
      </div>

      <!-- This Week -->
      <div class="p-4 border rounded-xl">
        <h3 class="font-semibold mb-2">This Week Todos</h3>

        <ul class="space-y-2">
          <li v-for="t in dashboard.thisWeek" :key="t.title">• {{ t.title }}</li>
        </ul>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { useDashboardStore } from "../stores/dashboard";
import AppLayout from "../layouts/AppLayout.vue";

const dashboard = useDashboardStore();

onMounted(() => {
  dashboard.fetchDashboard();
});
</script>

<style scoped>
.card {
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 12px;
}
</style>
