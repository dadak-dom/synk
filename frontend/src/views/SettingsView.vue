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
  <main>
    <div class="title">
      <h1>Settings</h1>
    </div>
    <div class="settings">
      <p>My private IP: {{ myPrivateIP }}</p>
    </div>
  </main>
</template>

<style scoped>
.settings {
  margin-top: 20px;
}
</style>
