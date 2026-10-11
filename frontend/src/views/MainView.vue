<script lang="ts" setup>
import { RunSynkOnPeer, GetConfigValueString, SetConfigItemString, SetConfigItemStringList, GetConfigValueStringList, GetPeerList } from "../../wailsjs/go/main/App";
import { inject, onMounted, ref } from "vue";
import { RouterLink } from "vue-router";
import { RemotePeer } from "../types/remotePeer";

// Theme piping
const props = defineProps([
  "firstColor",
  "secondColor",
  "thirdColor",
  "fourthColor",
  "textColor",
  "imageFilter",
  "backgroundColor",
  "borderColor",
  "onlinePeers",
]);

// TODO: This is how I could get the file information from a remote peer.
async function synk() {
  console.log("TESTING: ", selectedPeers.value)
  const peers = selectedPeers.value;
  const my_device_name = await GetConfigValueString("device_name.txt")
  peers.forEach(async (p) => {
    let url = "http://" + p.IP + ":8080/getSharedFolder";
    let response = await fetch(url, {
      headers: {
        'DeviceName':my_device_name
      }
    });
    let sharedFolderContents = await response.json();
    console.log(sharedFolderContents);
    const success = await RunSynkOnPeer(
      "http://" + p.IP + ":8080",
      sharedFolderContents,
    );
  });
}
const selectedPeers = ref<RemotePeer[]>([]);

// const peers = inject<RemotePeer[]>("peers");
// const peers = ref<RemotePeer[]>([])

const enableAutoSynk = ref<boolean>(false);

const deviceName = inject<string>("deviceName");

function handleAutoSynk() {
  if(enableAutoSynk.value == false) {
    SetConfigItemString("enable_auto_synk.txt", "false")
  } else {
    SetConfigItemString("enable_auto_synk.txt", "true")
  }
}

// Trusted network check
function checkIfNetworkTrusted() {
  const cn = ref<string>(inject("networkName") ?? "");
  const tn = ref<Array<string>>(inject("trustedNetworks") ?? []);
  console.log("cn: ", cn, "tn: ", tn);
  if (tn.value.includes(cn.value)) {
    console.log("network trusted");
    return true;
  }
  console.log("not trusted");
  return false;
}

// Update trusted peers config item
function updateTrustedPeers(){
  SetConfigItemStringList("trusted_devices.jsonl", selectedPeers.value.map((peer) => JSON.stringify(peer)))
  console.log("Selected peers: ", selectedPeers.value)
  console.log("Function in question: " , selectedPeers.value.filter((function (p: RemotePeer) {
              return props.onlinePeers.indexOf(p) >= 0
            })))
}

onMounted(() => {
  // check if autoSynk enabled or not
  GetConfigValueString("enable_auto_synk.txt").then((value) => enableAutoSynk.value = (value === "true"))
  // get list of trusted devices
  GetConfigValueStringList("trusted_devices.jsonl").then((value) => { selectedPeers.value = value.map((p) => JSON.parse(p))})
})
</script>

<template>
  <main>
    <div class="main-view-wrapper">
      <h1 id="logo" class="title">Synk</h1>
      <div class="peer-list-wrapper">
        <p v-if="onlinePeers == null">Scanning for peers...</p>
        <p v-else-if="onlinePeers !== null && onlinePeers.length == 0">
          No peers found.
          <RouterLink to="/settings"
            ><span>Check your connection.</span></RouterLink
          >
        </p>
        <input
          v-model="selectedPeers"
          v-for="peer in onlinePeers"
          type="checkbox"
          :value="peer"
          @change="updateTrustedPeers"
        />
        <label v-for="peer in onlinePeers" :for="peer.DeviceName">{{ peer.DeviceName }}</label>
        <div v-if="onlinePeers != undefined && selectedPeers.length > onlinePeers.length">
          <div>Spacer: offline peers</div>
          <input
            v-model="selectedPeers"
            v-for="peer in selectedPeers.filter(function (p) {
              return onlinePeers ? onlinePeers.indexOf(p) < 0 : true
            })"
            type="checkbox"
            :value="peer"
            @change="updateTrustedPeers"
          />
          <label v-for="peer in selectedPeers" :for="peer.DeviceName">{{ peer.DeviceName }}</label>
        </div>
      </div>

      <div class="option">
          <div>Auto-Synk</div>
          <label class="switch">
            <input
              type="checkbox"
              v-bind:checked="enableAutoSynk"
              v-model="enableAutoSynk"
              @change="handleAutoSynk"
            />
            <span class="slider round"></span>
          </label>
        </div>
      <button
        class="synk-button"
        v-if="
          onlinePeers != undefined &&
          selectedPeers.some(p =>
            onlinePeers.some((peer : RemotePeer) => peer.IP === p.IP)
          ) &&
          checkIfNetworkTrusted() &&
          !enableAutoSynk
        "
      >
        <img
          id="main-synk-button"
          @click="synk"
          src="../assets/images/refresh.png"
        />
      </button>
      <button class="synk-button-disabled" v-else>
        <img
          id="main-synk-button-disabled"
          @click=""
          src="../assets/images/refresh.png"
        />
      </button>
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
    v-bind(firstColor) 0,
    v-bind(secondColor) 20%,
    v-bind(thirdColor) 40%,
    v-bind(fourthColor) 100%
  );
  margin: auto;
  width: 200px;
  height: 200px;
  border-radius: 200px;
  display: flex;
  box-shadow: 0 0 10px 1px v-bind(borderColor);
}
.synk-button-disabled {
  margin-top: 200px;
  background: linear-gradient(
    200deg,
    v-bind(firstColor) 0,
    v-bind(secondColor) 20%,
    v-bind(thirdColor) 40%,
    v-bind(fourthColor) 100%
  );
  margin: auto;
  width: 200px;
  height: 200px;
  border-radius: 200px;
  border: 1px solid v-bind(borderColor);
  display: flex;
  box-shadow: 0 0 10px 1px v-bind(borderColor);
}

.peer-list-wrapper {
  border: solid 1px v-bind(borderColor);

  box-shadow:
    0 0 5px 1px v-bind(borderColor),
    0 0 50px 10px v-bind(firstColor);
  border-radius: 20px;
  width: 80%;
  margin: auto auto 60px auto;
  background-color: v-bind(backgroundColor);
}

#main-synk-button {
  width: 128px;
  margin: auto;
  filter: v-bind(imageFilter);
}

#main-synk-button-disabled {
  width: 128px;
  margin: auto;
  filter: v-bind(imageFilter);
  opacity: 0.5;
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
    200deg,
    v-bind(firstColor) 0,
    v-bind(secondColor) 20%,
    v-bind(thirdColor) 80%,
    v-bind(fourthColor) 100%
  );
  margin: 40px auto;
  width: 50%;
  padding: 20px 0;
  border-radius: 40px;
  border: 1px solid v-bind(borderColor);
  box-shadow: 0 0 20px 1px v-bind(borderColor);
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
