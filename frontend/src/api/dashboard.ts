import { api } from "../lib/axios";

export const getSummary = () => api.get("/dashboard/summary");

export const getAISummary = () => api.get("/dashboard/summary/ai");

export const getThisWeekTodos = () => api.get("/dashboard/this-week-todos");

export const getDailyRecommendation = () =>
  api.get("/dashboard/daily-recommendation");

export const streamAISummary = async (
  onMessage: (data: string) => void,
  onError?: () => void,
  onEnd?: () => void,
) => {
  const token = localStorage.getItem("access_token");
  const baseURL = import.meta.env.VITE_API_URL;

  fetch(`${baseURL}/dashboard/summary/ai`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })
    .then((response) => {
      if (!response.ok) throw new Error("Failed to fetch AI summary");

      const reader = response.body?.getReader();
      if (!reader) throw new Error("No response body");

      const decoder = new TextDecoder();
      let buffer = "";

      const read = () => {
        reader.read().then(({ done, value }) => {
          if (done) {
            onEnd?.();
            return;
          }

          buffer += decoder.decode(value, { stream: true });

          const events = buffer.split("\n\n");

          buffer = events.pop() ?? "";

          for (const event of events) {
            const lines = event.split("\n");

            for (const line of lines) {
              if (!line.startsWith("data:")) continue;
              onMessage(line.slice(5).trim() + "\n");
            }
          }

          read();
        });
      };

      read();
    })
    .catch((e) => {
      console.error("Stream error:", e);
      onError?.();
    });
};
