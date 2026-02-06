<script lang="ts" setup>
import { prependToMemberExpression } from '@babel/types';
import { GetPeerList } from '../../wailsjs/go/main/App';
import { onMounted, ref } from 'vue';
const peers = ref<string[]>([])

async function updatePeerList() {
  const result = await GetPeerList()
  peers.value = result
  console.log("Peers: ", peers.value)
}
onMounted(() => {
  console.log("Mounted peerlist")
  setInterval(updatePeerList, 500)
})


</script>

<template>
  <main><p>Peer list goes here</p>
  <div v-for="peer in peers"><p>{{ peer }}</p></div>
  </main>
</template>

<style scoped></style>
