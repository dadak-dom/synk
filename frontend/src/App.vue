<script lang="ts" setup>
import { ref } from "vue";
import { RouterLink } from "vue-router";

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
  <router-view v-slot="{ Component }">
    <transition name="fade" mode="out-in">
      <component :is="Component" :key="$route.path"></component>
    </transition>
  </router-view>
</template>

<style>
main {
  display: flex;
  flex-direction: column;
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
}

nav .nav-item img {
  width: 32px;
  filter: invert();
}

nav #navbar-burger {
  width: 32px;
  filter: invert();
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
  /* background-color: green; */
}

.wrapper {
  /* max-width: 80vw; */
  height: 90%;
  margin: auto;
  display: flex;
  justify-content: space-evenly;
  /* flex-grow: 0; */
}

/* #logo {
  display: block;
  width: 50%;
  height: 50%;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
} */

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
