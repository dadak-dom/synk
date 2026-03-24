<script lang="ts" setup>
import { onMounted, provide, reactive, ref } from "vue";
import { RouterLink } from "vue-router";
import {
  GetConfigValueString,
  GetPeerList,
  GetTheme,
  SetConfigItemString,
} from "../wailsjs/go/main/App";

// Router / Navbar stuff

const openNavBar = ref<boolean>(false);
let closeTimeout: ReturnType<typeof setTimeout> | null = null;

function mouseIn() {
  if (closeTimeout) {
    clearTimeout(closeTimeout);
    closeTimeout = null;
  }
  setTimeout(() => {
    openNavBar.value = true;
  }, 100);
}
function mouseOut() {
  closeTimeout = setTimeout(() => {
    openNavBar.value = false;
    closeTimeout = null;
  }, 300);
}

// Peerlist tracking

const peers = ref<string[] | null>(null);
const selectedPeers = ref<string[]>([]);

provide("peers", peers);

async function updatePeerList() {
  const result = await GetPeerList();
  peers.value = result;
  console.log("Peers: ", peers.value, "Selected peers: ", selectedPeers.value);
}

// Config items in settings
const autoIgnoreAll = ref<boolean>();
provide("autoIgnoreAll", autoIgnoreAll);

// Themes
import { Theme, GetThemeInterface } from "./interfaces/theme";

const theme = ref<Theme>(GetThemeInterface("Dark"));

const themeName = ref<string>(""); // CSS variables
const firstColor = ref<string>("");
const secondColor = ref<string>("");
const thirdColor = ref<string>("");
const fourthColor = ref<string>("");
const textColor = ref<string>("");
const themeFilter = ref<string>("");

provide("theme", theme.value);

function updateTheme(th: string) {
  let new_theme = GetThemeInterface(th);
  // CSS variables to update
  firstColor.value = new_theme.firstColor;
  secondColor.value = new_theme.secondColor;
  thirdColor.value = new_theme.thirdColor;
  fourthColor.value = new_theme.fourthColor;
  textColor.value = new_theme.textColor;
  themeName.value = new_theme.name;
  themeFilter.value = new_theme.filter;

  // ref to update
  theme.value.name = new_theme.name;
  // Save the choice to config
  SetConfigItemString("theme.txt", th);
}

onMounted(() => {
  console.log("Mounting App.vue...");
  GetTheme().then((th) => {
    updateTheme(th);
  });
  GetConfigValueString("auto_ignore_all.txt").then((aia) =>
    aia == "true"
      ? (autoIgnoreAll.value = true)
      : (autoIgnoreAll.value = false),
  );
  setInterval(updatePeerList, 3000);
});
</script>

<template>
  <Transition name="slide-fade">
    <nav v-if="openNavBar" @mouseleave="mouseOut" @mouseenter="mouseIn">
      <RouterLink class="nav-item" to="/"
        ><img style="" src="./assets/images/home.png"
      /></RouterLink>
      <RouterLink class="nav-item" to="/folder"
        ><img src="./assets/images/nav_folder.png"
      /></RouterLink>
      <RouterLink class="nav-item" to="/about"
        ><img src="./assets/images/about.png"
      /></RouterLink>
      <RouterLink class="nav-item" to="/settings"
        ><img src="./assets/images/settings.png"
      /></RouterLink>
    </nav>
    <nav v-else @mouseenter="mouseIn">
      <img id="navbar-burger" src="./assets/images/navbar_icon.png" />
    </nav>
  </Transition>
  <router-view
    v-slot="{ Component }"
    @update-theme="updateTheme"
    :firstColor="firstColor"
    :secondColor="secondColor"
    :thirdColor="thirdColor"
    :fourthColor="fourthColor"
    :textColor="textColor"
    :imageFilter="themeFilter"
  >
    <transition name="fade" mode="out-in">
      <component :is="Component" :key="$route.path"></component>
    </transition>
  </router-view>
</template>

<style>
main {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: linear-gradient(
    200deg,
    v-bind(firstColor) 0,
    v-bind(secondColor) 20%,
    v-bind(thirdColor) 40%,
    v-bind(fourthColor) 100%
  );
  color: v-bind(textColor);
}

.outer-view {
  width: 100%;
  display: flex;
  justify-content: center;
  height: 100%;
}

nav {
  display: flex;
  justify-content: space-evenly;
  flex-direction: column;
  position: absolute;
  z-index: 100;
  gap: 10px;
  margin-top: 10px;
  margin-left: 10px;
  background: linear-gradient(
    200deg,
    v-bind(firstColor) 0,
    v-bind(secondColor) 20%,
    v-bind(thirdColor) 40%,
    v-bind(fourthColor) 100%
  );
  border: solid 1px grey;
  border-radius: 10px;
  padding: 5px;
}

nav .nav-item img {
  width: 32px;
  filter: v-bind(themeFilter);
}

nav #navbar-burger {
  width: 32px;
  filter: v-bind(themeFilter);
  cursor: pointer;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease-in-out;
}
.fade-enter,
.fade-leave-to {
  opacity: 0;
}

.main-app {
  display: flex;
  justify-content: space-around;
  height: 100%;
  max-width: 80%;
  margin: auto;
}

.wrapper {
  height: 90%;
  margin: auto;
  display: flex;
  justify-content: space-evenly;
}

.slide-fade-enter-active {
  transition: all 0.3s ease-out;
}

.slide-fade-leave-active {
  transition: all 0.3s ease-in-out;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateX(-60px);
  opacity: 0;
}
</style>
