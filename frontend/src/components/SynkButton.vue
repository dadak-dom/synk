<script lang="ts" setup>
import { GetPeerList } from "../../wailsjs/go/main/App";
import { RunSynkOnPeer } from "../../wailsjs/go/main/App";
import { TestLANDiscovery } from "../../wailsjs/go/main/App";
import { onMounted, ref } from 'vue';

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


const peers = ref<string[]>([])
const selectedPeers = ref<string[]>([])

async function updatePeerList() {
  const result = await GetPeerList()
  peers.value = result
  console.log("Peers: ", peers.value, "Selected peers: ", selectedPeers.value)

}
onMounted(() => {
  console.log("Mounted peerlist")
  setInterval(updatePeerList, 500)
})
</script>

<template>
  <main>
    <p>Peer list goes here</p>
    <input v-model="selectedPeers" v-for="peer in peers" type="checkbox" :name="peer" :value="peer"/>
    <label v-for="peer in peers" :for="peer">{{ peer }}</label>

    <div>
      <img @click="synk" src="../assets/images/sink.png" />
       <!-- <img @click="testLan" src="../assets/images/sink.png" /> -->
      <h1 id="app-name">"Synk"</h1>
    </div>
  </main>
</template>

<style scoped>
@font-face {
  font-family: "LogoFont";
  src: url("../assets/fonts/typewriter.otf") format("opentype");
}
img {
  height: 25%;
  width: 25%;
  margin: auto;
}
#app-name {
  font-family: "LogoFont";
  color: black;
}
</style>
