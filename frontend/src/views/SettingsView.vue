<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { GetLocalIP } from "../../wailsjs/go/main/App";
const myPrivateIP = ref<string>("");

onMounted(() => {
  GetLocalIP().then((ip) => {
    myPrivateIP.value = ip;
  });
});
</script>

<template>
  <div class="settings">
    <h1>Settings</h1>
    <p>My private IP: {{ myPrivateIP }}</p>
  </div>
</template>

<style scoped></style>
