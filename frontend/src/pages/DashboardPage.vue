<style scoped>
.card {
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 12px;
}
</style>

<template>
  <AppLayout>
    <div class="space-y-6">
      <div>
        <div class="text-2xl font-bold">{{ getGreeting() }}, user! 👋</div>
        <div class="text-gray-500 dark:text-gray-400">
          {{ dashboard.recommendation ?? "Here's your summary for today." }}
        </div>
      </div>
      <div class="flex flex-col md:flex-row justify-start w-full gap-4">
        <div class="flex flex-col gap-4 w-full md:w-2/3">
          <!-- Summary Cards -->
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <SummaryCard title="Total" :value="dashboard.summary?.total ?? 0" />
            <SummaryCard
              title="Completed"
              :value="dashboard.summary?.completed_count ?? 0"
            />
            <SummaryCard
              title="Pending"
              :value="dashboard.summary?.pending_count ?? 0"
            />
            <SummaryCard
              title="Completion Rate"
              :value="`${dashboard.summary?.completion_rate ?? 0}%`"
            />
          </div>
          <!-- Summary -->
          <UCard>
            <div class="flex justify-start items-center space-x-2 gap-2">
              <UButton
                class="flex items-center gap-2 mb-2 disabled:opacity-50 disabled:cursor-not-allowed"
                @click="fetchAISummary"
                :disabled="isLoadingSummarize || isSummaryVisible"
              >
                {{
                  isSummaryVisible
                    ? "Todos summarized!"
                    : isLoadingSummarize
                      ? "Loading..."
                      : "Show all todos summary..."
                }}
                <Astroid />
              </UButton>
            </div>
            <p class="text-sm text-gray-500 dark:text-gray-400">
              <span v-html="dashboard.aiSummary" />
            </p>
          </UCard>
        </div>

        <!-- This Week -->
        <div
          class="flex flex-col gap-4 w-full md:w-1/3 md:ml-4"
          v-if="dashboard.thisWeek"
        >
          <UCard>
            <template #header>
              <h3 class="text-md font-semibold text-gray-900 dark:text-white">
                Todos This Week
              </h3>
            </template>
            <div v-for="t in dashboard.thisWeek.days" :key="t.date">
              <p class="font-medium">
                {{ new Date(t.date).toLocaleDateString() }}
              </p>
              <ul
                v-for="todo in t.todos"
                :key="todo.id"
                class="list-disc list-inside"
              >
                <li>{{ todo.title }}</li>
              </ul>
            </div>
          </UCard>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useDashboardStore } from "../stores/dashboard";
import { streamAISummary } from "../api/dashboard";
import AppLayout from "../layouts/AppLayout.vue";
import { Astroid } from "@lucide/vue";
import SummaryCard from "../components/dashboard/SummaryCard.vue";

const dashboard = useDashboardStore();
const isLoadingSummarize = ref(false);
const isSummaryVisible = ref(false);

const fetchAISummary = async () => {
  if (isLoadingSummarize.value) return;

  isLoadingSummarize.value = true;
  let fullMessage = "";

  streamAISummary(
    (data) => {
      fullMessage += data;
      dashboard.aiSummary = fullMessage;
    },
    () => {
      isLoadingSummarize.value = false;
    },
    () => {
      isLoadingSummarize.value = false;
      isSummaryVisible.value = true;
    },
  );
};

const getGreeting = () => {
  const hour = new Date().getHours();
  if (hour < 12) return "Good morning";
  if (hour < 18) return "Good afternoon";
  return "Good evening";
};

onMounted(() => {
  dashboard.fetchDashboard();
});
</script>
