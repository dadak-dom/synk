<script lang="ts" setup>
import { GetPeerList } from "../../wailsjs/go/main/App";
import { RunSynkOnPeer } from "../../wailsjs/go/main/App";
import { onMounted, ref } from "vue";

// TODO: This is how I could get the file information from a remote peer.
async function synk() {
  const peers = selectedPeers.value;
  peers.forEach(async (p) => {
    let url = "http://" + p + ":8080/getSharedFolder";
    let response = await fetch(url);
    let sharedFolderContents = await response.json();
    console.log(sharedFolderContents);
    // return
    RunSynkOnPeer("http://" + p + ":8080", sharedFolderContents);
  });
}

const peers = ref<string[]>([]);
const selectedPeers = ref<string[]>([]);

async function updatePeerList() {
  const result = await GetPeerList();
  peers.value = result;
  console.log("Peers: ", peers.value, "Selected peers: ", selectedPeers.value);
}

onMounted(() => {
  console.log("Mounted peerlist");
  setInterval(updatePeerList, 500);
});
</script>

<template>
  <main>
    <div class="main-view-wrapper">
      <h1 id="logo">Synk</h1>
      <p>Peer list goes here</p>
      <input
        v-model="selectedPeers"
        v-for="peer in peers"
        type="checkbox"
        :name="peer"
        :value="peer"
      />
      <label v-for="peer in peers" :for="peer">{{ peer }}</label>

      <div class="synk-button">
        <img
          id="main-synk-button"
          @click="synk"
          src="../assets/images/refresh.png"
        />
        <!-- <h1 id="app-name">"Synk"</h1> -->
      </div>
    </div>
  </main>
</template>

<style scoped>
@font-face {
  font-family: "LogoFont";
  src: url("../assets/fonts/typewriter.otf") format("opentype");
}

@font-face {
  font-family: "FrutigerAero";
  src: url("../assets/fonts/AerobicsRegular.ttf") format("truetype");
}

.synk-button {
  cursor: pointer;
  margin-top: 200px;
}

#main-synk-button {
  height: 25%;
  width: 25%;
  margin: auto;
  filter: contrast(100) invert();

  /* rotate: 0deg; */
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
    183deg,
    rgba(148, 148, 148, 0.9) 0,
    rgba(7, 7, 7, 0.93) 20%,
    rgba(19, 19, 19, 0.9) 40%,
    rgba(30, 30, 30, 0.93) 100%
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
