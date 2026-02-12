<script lang="ts" setup>
import { GetPeerList } from "../../wailsjs/go/main/App";
import { RunSynkOnPeer } from "../../wailsjs/go/main/App";
import { onMounted, ref } from "vue";
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
    RunSynkOnPeer("http://" + p + ":8080", sharedFolderContents);
  });
}

const peers = ref<string[] | null>(null);
const selectedPeers = ref<string[]>([]);

async function updatePeerList() {
  const result = await GetPeerList();
  peers.value = result;
  console.log("Peers: ", peers.value, "Selected peers: ", selectedPeers.value);
}

onMounted(() => {
  console.log("Mounted peerlist");
  setInterval(updatePeerList, 3000);
});
</script>

<template>
  <main>
    <div class="main-view-wrapper">
      <h1 id="logo" class="title">Synk</h1>
      <div class="peer-list-wrapper">
        <p v-if="peers == null">Scanning for peers...</p>
        <p v-else-if="peers.length == 0">
          No peers found.
          <RouterLink to="/settings"
            ><span
              style="
                border-bottom: dotted 1px gray;
                color: white;
                text-decoration: none;
              "
              >Check your connection.</span
            ></RouterLink
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

      <div class="synk-button">
        <img
          id="main-synk-button"
          @click="synk"
          src="../assets/images/refresh.png"
        />
      </div>
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
    rgba(148, 148, 148, 0.6) 0,
    rgba(7, 7, 7, 0.6) 20%,
    rgba(19, 19, 19, 0.6) 40%,
    rgba(105, 102, 102, 0.6) 100%
  );
  margin: auto;
  width: 200px;
  height: 200px;
  border-radius: 200px;
  display: flex;
}

.peer-list-wrapper {
  border: solid 1px darkgray;
  width: 80%;
  margin: auto auto 60px auto;
  background-color: rgba(40, 40, 40, 0.5);
}

#main-synk-button {
  /* height: 25%; */
  /* width: 25%; */
  width: 128px;
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
