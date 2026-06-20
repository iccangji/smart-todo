import { api } from "../lib/axios";

export const getSummary = () =>
  api.get("/dashboard/summary");

export const getAISummary = () =>
  api.get("/dashboard/summary/ai");

export const getThisWeekTodos = () =>
  api.get("/dashboard/this-week-todos");

export const getDailyRecommendation = () =>
  api.get("/dashboard/daily-recommendation");
