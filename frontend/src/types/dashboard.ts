import type { Todo } from "./todo";

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

export interface ThisWeekTodo {
  days: DailyTodo[];
}

interface DailyTodo {
  date: string;
  todos: Todo[];
}
