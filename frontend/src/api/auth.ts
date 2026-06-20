import { api } from "../lib/axios";

export const login = (payload: {
  email: string;
  password: string;
}) => {
  return api.post("/auth/login", payload);
};

export const register = (payload: {
  email: string;
  password: string;
  name: string;
}) => {
  return api.post("/auth/register", payload);
};

export const refresh = () => {
  return api.post("/auth/refresh");
};
