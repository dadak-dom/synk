<script lang="ts" setup>
import { inject, onMounted, ref } from "vue";
import {
  GetLocalIP,
  GetCurrentNetworkName,
  GetConfigValueStringList,
  SetConfigItemStringList,
} from "../../wailsjs/go/main/App";
import { Theme, AllThemeNames } from "../interfaces/theme";
const myPrivateIP = ref<string>("");

const userTheme = inject<Theme>("theme");
const emit = defineEmits(["updateTheme"]);
const selectedTheme = ref<string>(userTheme?.name ? userTheme.name : "");

const currentNetwork = ref<string>(inject("networkName") ?? "");
const trustedNetworks = ref<Array<string>>(inject("trustedNetworks") ?? []);

function themeInput(th: string) {
  emit("updateTheme", th);
  console.log("change event", th);
}

function addTrustedNetwork(network: string) {
  GetConfigValueStringList("trusted_networks.jsonl").then(
    (networks: Array<string>) => {
      SetConfigItemStringList(
        "trusted_networks.jsonl",
        networks.concat([network]),
      );
      trustedNetworks.value = networks.concat([network]);
    },
  );
}

function resetTrustedNetworks() {
  SetConfigItemStringList("trusted_networks.jsonl", []).then(
    () => (trustedNetworks.value = []),
  );
}

function checkIfAlreadyTrusted() {
  console.log(
    "TESTING!!!",
    currentNetwork.value,
    trustedNetworks.value,
    trustedNetworks.value.includes(currentNetwork.value),
  );
  if (trustedNetworks.value.includes(currentNetwork.value)) {
    return true;
  }
  return false;
}

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
      <p>Current network: {{ inject("networkName") }}</p>
      <p>Trusted networks list: {{ trustedNetworks }}</p>
      <button
        @click="addTrustedNetwork(currentNetwork)"
        :disabled="checkIfAlreadyTrusted()"
      >
        Trust my network
      </button>
      <button @click="resetTrustedNetworks">Reset trusted networks</button>
      <div class="theme-selector">
        <label for="themes">Theme:</label>
        <select
          name="themes"
          v-model="selectedTheme"
          @change="themeInput(selectedTheme)"
        >
          <template v-for="theme in AllThemeNames">
            <option :value="theme">{{ theme }}</option>
            <p>{{ userTheme?.name }}</p>
          </template>
        </select>
      </div>
    </div>
  </main>
</template>

<style scoped>
.settings {
  margin-top: 20px;
}
.theme-selector {
  display: flex;
  justify-content: center;
  gap: 10px;
}
</style>
