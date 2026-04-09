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
const revealIP = ref<boolean>(false);

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

// theme
const props = defineProps([
  "firstColor",
  "secondColor",
  "thirdColor",
  "fourthColor",
  "textColor",
  "imageFilter",
  "backgroundColor",
  "borderColor",
]);

onMounted(() => {
  GetLocalIP().then((ip) => {
    myPrivateIP.value = ip;
  });
});
</script>

<template>
  <main>
    <div class="title">
      <h1>Settings<br /></h1>
    </div>
    <div class="settings">
      <div class="private-ip-wrapper">
        <div>My private IP:</div>
        <div
          @mouseenter="revealIP = true"
          @mouseleave="revealIP = false"
          class="private-ip"
        >
          {{ myPrivateIP }}
          <div :hidden="revealIP" class="censor-text">.</div>
        </div>
      </div>
      <div class="trusted-networks-wrapper">
        <div>Current network: {{ inject("networkName") }}</div>
        <div>Trusted networks list: {{ trustedNetworks }}</div>
        <button
          @click="addTrustedNetwork(currentNetwork)"
          :disabled="checkIfAlreadyTrusted()"
        >
          Trust my network
        </button>
        <button
          @click="resetTrustedNetworks"
          :disabled="trustedNetworks.length == 0"
        >
          Reset trusted networks
        </button>
      </div>
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
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 10px;
}
.theme-selector {
  display: flex;
  justify-content: center;
  gap: 10px;
}
.trusted-networks-wrapper {
}
.private-ip-wrapper {
  display: flex;
  flex-direction: row;
  height: fit-content;
  width: fit-content;
  gap: 10px;
}
.private-ip {
  position: relative;
  width: fit-content;
  cursor: default;
}
.censor-text {
  background-color: v-bind(textColor);
  position: absolute;
  top: 0;
  width: 100%;
  height: 100%;
}
</style>
