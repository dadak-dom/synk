<script lang="ts" setup>
import FolderSelector from "./components/FolderSelector.vue";
import PeerList from "./components/PeerList.vue";
import SynkButton from "./components/SynkButton.vue";
</script>

<template>
  <!-- <img id="logo" alt="Wails logo" src="./assets/images/logo-universal.png" /> -->
  <!-- <AppLogo /> -->
  <!-- <SetDirectory /> -->
  <!-- <SynkButton /> -->
  <div class="main-app">
    <div class="wrapper">
      <FolderSelector />
      <PeerList />
      <SynkButton />
    </div>
  </div>
</template>

<style>
* {
  border: solid 1px red;
}

main {
  width: 33.33%;
  max-height: 580px;
  margin: auto;
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
