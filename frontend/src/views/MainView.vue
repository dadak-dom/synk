<script lang="ts" setup>
import { RunSynkOnPeer } from "../../wailsjs/go/main/App";
import { inject, ref } from "vue";
import { RouterLink } from "vue-router";

// TODO: This is how I could get the file information from a remote peer.
async function synk() {
  const peers = selectedPeers.value;
  peers.forEach(async (p) => {
    let url = "http://" + p + ":8080/getSharedFolder";
    let response = await fetch(url);
    let sharedFolderContents = await response.json();
    console.log(sharedFolderContents);
    // return
    const success = await RunSynkOnPeer(
      "http://" + p + ":8080",
      sharedFolderContents,
    );
    if (success) {
      alert("Synk with " + p + " successful!");
    } else {
      alert("Synk with " + p + " failed.");
    }
  });
}
const selectedPeers = ref<string[]>([]);

const peers = inject<string[] | null>("peers");

// Theme piping
defineProps([
  "firstColor",
  "secondColor",
  "thirdColor",
  "fourthColor",
  "textColor",
  "imageFilter",
]);
</script>

<template>
  <main>
    <div class="main-view-wrapper">
      <h1 id="logo" class="title">Synk</h1>
      <div class="peer-list-wrapper">
        <p v-if="peers == null">Scanning for peers...</p>
        <p v-else-if="peers !== null && peers.length == 0">
          No peers found.
          <RouterLink to="/settings"
            ><span>Check your connection.</span></RouterLink
          >
        </p>
        <input
          v-model="selectedPeers"
          v-for="peer in peers"
          type="checkbox"
          :name="peer"
          :value="peer"
        />
        <label v-for="peer in peers" :for="peer">{{ peer }}</label>
      </div>

      <button class="synk-button" v-if="selectedPeers.length > 0">
        <img
          id="main-synk-button"
          @click="synk"
          src="../assets/images/refresh.png"
        />
      </button>
      <button class="synk-button-disabled" v-else>
        <img
          id="main-synk-button-disabled"
          @click="synk"
          src="../assets/images/refresh.png"
        />
      </button>
    </div>
  </main>
</template>

<style scoped>
@font-face {
  font-family: "LogoFont";
  src: url("../assets/fonts/typewriter.otf") format("opentype");
}

.synk-button {
  cursor: pointer;
  margin-top: 200px;
  background: linear-gradient(
    180deg,
    v-bind(firstColor) 0,
    v-bind(secondColor) 20%,
    v-bind(thirdColor) 40%,
    v-bind(fourthColor) 100%
  );
  margin: auto;
  width: 200px;
  height: 200px;
  border-radius: 200px;
  display: flex;
}
.synk-button-disabled {
  margin-top: 200px;
  background: linear-gradient(
    200deg,
    v-bind(firstColor) 0,
    v-bind(secondColor) 20%,
    v-bind(thirdColor) 40%,
    v-bind(fourthColor) 100%
  );
  margin: auto;
  width: 200px;
  height: 200px;
  border-radius: 200px;
  border: 1px solid black;
  display: flex;
}

.peer-list-wrapper {
  border: solid 1px darkgray;
  width: 80%;
  margin: auto auto 60px auto;
  background-color: v-bind(firstColor);
}

#main-synk-button {
  width: 128px;
  margin: auto;
  filter: v-bind(imageFilter);
}

#main-synk-button-disabled {
  width: 128px;
  margin: auto;
  filter: contrast(100) invert();
  opacity: 0.5;
}

.main-view-wrapper {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

#logo {
  font-family: "FrutigerAero";
  background: linear-gradient(
    180deg,
    v-bind(firstColor) 0,
    v-bind(secondColor) 20%,
    v-bind(thirdColor) 40%,
    v-bind(fourthColor) 100%
  );
  margin: 40px auto;
  width: 50%;
  padding: 20px 0;
  border-radius: 40px;
}

#main-synk-button:hover {
  transform: rotate(360deg);
  animation: spin 0.5s ease-in-out;
  animation-iteration-count: infinite;
}
#app-name {
  font-family: "LogoFont";
  color: black;
}

@keyframes spin {
  from {
    rotate: 0deg;
  }
  to {
    rotate: 360deg;
  }
}
</style>
