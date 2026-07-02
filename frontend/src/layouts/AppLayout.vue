<script setup lang="ts">
import { computed, ref } from "vue";
import type { NavigationMenuItem, SidebarProps } from "@nuxt/ui";
import { useRoute } from "vue-router";

const open = ref(false);

const route = useRoute();
defineProps<Pick<SidebarProps, "variant" | "collapsible" | "side">>();
const toggle = () => {
  open.value = !open.value;
};
const items = computed<NavigationMenuItem[]>(() => [
  {
    label: "Dashboard",
    to: "/",
    icon: "i-lucide-home",
    active: route.path === "/",
  },
  {
    label: "Todos",
    to: "/todos",
    icon: "i-lucide-check-square",
    active: route.path === "/todos",
  },
]);
</script>

<template>
  <div class="flex flex-col flex-1">
    <UHeader title="Smart Todo" toggle-side="left" :ui="{ container: 'px-4!' }">
      <template #toggle>
        <UButton
          icon="i-lucide-panel-left"
          color="neutral"
          variant="ghost"
          aria-label="Toggle sidebar"
          @click="toggle"
        />
      </template>
      <template #right>
        <UColorModeButton />
      </template>
    </UHeader>

    <div class="flex flex-1 min-h-0">
      <USidebar
        v-model:open="open"
        collapsible="icon"
        :ui="{
          gap: 'h-[calc(100%-var(--ui-header-height))]',
          container:
            'absolute top-(--ui-header-height) bottom-0 h-[calc(100%-var(--ui-header-height))]',
        }"
      >
        <UNavigationMenu
          :items="items"
          orientation="vertical"
          :ui="{ link: 'p-1.5 overflow-hidden' }"
        />
      </USidebar>

      <div class="flex-1 p-4">
        <div class="size-full">
          <slot />
        </div>
      </div>
    </div>
  </div>
</template>
