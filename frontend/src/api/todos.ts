import { api } from "../lib/axios";

export const getTodos = () => api.get("/todos");

export const getTodo = (id: string) => api.get(`/todos/${id}`);

export const createTodo = (payload: unknown) => api.post("/todos", payload);

export const updateTodo = (id: string, payload: unknown) => api.put(`/todos/${id}`, payload);

export const deleteTodo = (id: string) => api.delete(`/todos/${id}`);

export const breakdownTodo = (id: string) => api.get(`/todos/${id}/breakdown`);
