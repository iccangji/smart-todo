export interface DashboardSummary {
  total: number;
  completed_count: number;
  pending_count: number;

  pending_priority_count: {
    low: number;
    medium: number;
    high: number;
    urgent: number;
  };

  completion_rate: number;
  completed_today: number;
  completed_this_week: number;
}

export interface DailyRecommendation {
  message: string;
}
