<template>
  <div class="layout">
    <aside class="sidebar">
      <RouterLink to="/dashboard">Dashboard</RouterLink>
      <RouterLink to="/todos">Todos</RouterLink>
    </aside>

    <main class="content">
      <slot />
    </main>
  </div>
</template>

<style scoped>
.layout {
  display: flex;
  min-height: 100vh;
  padding: 0;
  width: 100%;
}

.sidebar {
  width: 240px;
  padding: 1rem;
  border-right: 1px solid #eee;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.content {
  flex: 1;
  padding: 2rem;
}
</style>
