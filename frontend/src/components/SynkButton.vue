<script lang="ts" setup>
import { GetPeerList } from "../../wailsjs/go/main/App";
import { RunSynkOnPeer } from "../../wailsjs/go/main/App";
import { TestLANDiscovery } from "../../wailsjs/go/main/App";

// TODO: This is how I could get the file information from a remote peer.
async function synk() {
  const peers = await GetPeerList();
  peers.forEach(async (p) => {
    let url = "http://" + p + "/getSharedFolder";
    let response = await fetch(url);
    let sharedFolderContents = await response.json();
    console.log(sharedFolderContents);
    return 
    RunSynkOnPeer("http://" + p, sharedFolderContents);
  });
}

function testLan() {
  TestLANDiscovery()
}
</script>

<template>
  <main>
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
