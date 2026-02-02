<script lang="ts" setup>
import AppLogo from "./components/AppLogo.vue";
import FolderSelector from "./components/FolderSelector.vue";
import PeerList from "./components/PeerList.vue";
</script>

<template>
  <!-- <img id="logo" alt="Wails logo" src="./assets/images/logo-universal.png" /> -->
  <AppLogo />
  <!-- <SetDirectory /> -->
  <!-- <SynkButton /> -->
  <div class="main-app">
    <div class="wrapper">
      <FolderSelector />
      <PeerList />
    </div>
  </div>
</template>

<style>
.main-app {
  display: flex;
  justify-content: space-around;
}

.wrapper {
  max-width: 80vw;
  display: flex;
  justify-content: space-around;
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
