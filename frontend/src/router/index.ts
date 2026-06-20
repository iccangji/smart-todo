import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      redirect: "/dashboard",
    },
    {
      path: "/dashboard",
      component: () => import("../pages/DashboardPage.vue"),
    },
    {
      path: "/todos",
      component: () => import("../pages/TodosPage.vue"),
    },
  ],
});

export default router;
