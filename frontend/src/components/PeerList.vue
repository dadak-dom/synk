<script lang="ts" setup>
import { GetPeerList } from "../../wailsjs/go/main/App";
import { onMounted, ref } from "vue";
import { RemotePeer } from "../types/remotePeer";

const peers = ref<RemotePeer[]>([]);

async function updatePeerList() {
  const result = await GetPeerList();
  peers.value = result;
  console.log("Peers: ", peers.value);
}
onMounted(() => {
  console.log("Mounted peerlist");
  setInterval(updatePeerList, 500);
});
</script>

<template>
  <main></main>
</template>

<style scoped></style>
