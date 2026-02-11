<script lang="ts" setup>
import { createMemoryHistory, createRouter, RouterLink } from "vue-router";
</script>

<template>
  <nav>
    <RouterLink to="/">Go Home</RouterLink>
    <RouterLink to="/folder"
      ><img src="./assets/images/folder.png"
    /></RouterLink>
    <RouterLink to="/about">Go to About</RouterLink>
    <RouterLink to="/settings">settings</RouterLink>
  </nav>
  <!-- <RouterView /> -->
  <router-view v-slot="{ Component }">
    <transition name="fade" mode="out-in">
      <component :is="Component" :key="$route.path"></component>
    </transition>
  </router-view>
</template>

<style>
* {
  /* border: solid 1px red; */
}

nav {
  /* height: 10%; */
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
  /* max-width: 80vw; */
  height: 90%;
  margin: auto;
  display: flex;
  justify-content: space-evenly;
  /* flex-grow: 0; */
}

#logo {
  display: block;
  width: 50%;
  height: 50%;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
}
</style>
