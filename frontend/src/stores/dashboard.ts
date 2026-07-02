import { defineStore } from "pinia";
import {
  getSummary,
  getDailyRecommendation,
  getThisWeekTodos,
  // getAISummary,
} from "../api/dashboard";
import type { DashboardSummary, ThisWeekTodo } from "../types/dashboard";

export const useDashboardStore = defineStore("dashboard", {
  state: () => ({
    summary: null as DashboardSummary | null,
    aiSummary: null as string | null,
    thisWeek: null as ThisWeekTodo | null,
    recommendation: null as string | null,
    loading: false,
  }),

  actions: {
    async fetchDashboard() {
      this.loading = true;

      try {
        const [summary, recommendation, week] = await Promise.all([
          getSummary(),
          getDailyRecommendation(),
          getThisWeekTodos(),
          // getAISummary(),
        ]);

        // console.log("Summary:", week.data.data.days);

        this.summary = summary.data.data;
        this.recommendation = recommendation.data.data.message;
        this.thisWeek = week.data.data;
        // this.aiSummary = ai.data.data;
      } finally {
        this.loading = false;
      }
    },
  },
});
