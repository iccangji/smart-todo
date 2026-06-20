import { defineStore } from "pinia";
import {
  getSummary,
  // getDailyRecommendation,
  getThisWeekTodos,
  // getAISummary,
} from "../api/dashboard";

export const useDashboardStore = defineStore("dashboard", {
  state: () => ({
    summary: null as any,
    recommendation: null as any,
    thisWeek: [] as any[],
    aiSummary: null as any,

    loading: false,
  }),

  actions: {
    async fetchDashboard() {
      this.loading = true;

      try {
        const [summary, week] = await Promise.all([
          getSummary(),
          // getDailyRecommendation(),
          getThisWeekTodos(),
          // getAISummary(),
        ]);

        console.log("Summary:", summary.data.data);

        this.summary = summary.data.data;
        // this.recommendation = rec.data.data;
        this.thisWeek = week.data.data;
        // this.aiSummary = ai.data.data;
      } finally {
        this.loading = false;
      }
    },
  },
});
