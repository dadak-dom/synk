<script lang="ts" setup>
import { GetPeerList } from "../../wailsjs/go/main/App";
import { onMounted, ref } from "vue";
const peers = ref<string[]>([]);

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
