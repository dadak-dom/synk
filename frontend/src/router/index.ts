import { createRouter, createWebHistory } from "vue-router";
import MainView from "../views/MainView.vue";
import AboutView from "../views/AboutView.vue";
import FolderSelectorView from "../views/FolderSelectorView.vue";
import SettingsView from "../views/SettingsView.vue";

const routes = [
  { path: "/", component: MainView },
  { path: "/about", component: AboutView },
  { path: "/folder", component: FolderSelectorView },
  { path: "/settings", component: SettingsView },
];

export const router = createRouter({
  history: createWebHistory(),
  routes,
});
