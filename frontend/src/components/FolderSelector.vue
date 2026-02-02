<script lang="ts" setup>
import { onMounted, ref } from "vue";
// import what I need from the backend...
import { FolderSelectorControl } from "../../wailsjs/go/main/App";

const enableSelector = ref(false);
const currentDir = ref();
const filesInCurrentDir = ref<string[]>([]);

const FolderSelectorCommands = {
  MOVE_UP: 0,
  MOVE_DOWN: 1,
  GO_HOME: 2,
  INIT: 3,
  SELECT: 4,
};

function moveUpDir() {
  FolderSelectorControl(
    currentDir.value,
    FolderSelectorCommands.MOVE_UP,
    "",
  ).then((value) => {
    currentDir.value = value.Directory;
    filesInCurrentDir.value = value.Files;
    console.log(filesInCurrentDir.value);
  });
}

function moveDownDir(f: string) {
  FolderSelectorControl(
    currentDir.value,
    FolderSelectorCommands.MOVE_DOWN,
    f,
  ).then((value) => {
    currentDir.value = value.Directory;
    filesInCurrentDir.value = value.Files;
    console.log(filesInCurrentDir.value);
  });
}

function goHome() {
  FolderSelectorControl("", FolderSelectorCommands.GO_HOME, "").then(
    (value) => {
      currentDir.value = value.Directory;
      filesInCurrentDir.value = value.Files;
      console.log(filesInCurrentDir.value);
    },
  );
}

function selectFolder() {
  enableSelector.value = false;
  FolderSelectorControl(currentDir.value, FolderSelectorCommands.SELECT, "");
}

function chooseNewDir() {
  enableSelector.value = true;
}

onMounted(() => {
  FolderSelectorControl("", FolderSelectorCommands.INIT, "").then((value) => {
    currentDir.value = value.Directory;
    filesInCurrentDir.value = value.Files;
    console.log(filesInCurrentDir.value);
    if (currentDir.value == "") {
      enableSelector.value = true;
    }
  });
});
</script>

<template>
  <main class="outer-div" v-if="enableSelector">
    <div class="nav-button-group">
      <p class="current-directory">Current directory: {{ currentDir }}</p>
      <button class="nav-button" @click="moveUpDir">↑ Move Up</button>
      <button class="nav-button" @click="goHome">Go Home</button>
      <button class="nav-button" @click="selectFolder">
        Select this folder
      </button>
    </div>

    <div
      class="folder-selection-box"
      v-if="filesInCurrentDir == undefined || filesInCurrentDir.length > 0"
    >
      <div v-for="file in filesInCurrentDir">
        <!-- <FolderSelectorItem :fileName="file" /> -->

        <button class="folder-button" @click="moveDownDir(file)">
          {{ file }}
          <img class="folder-image" src="../assets/images/folder.png" />
        </button>
      </div>
    </div>
    <div v-else class="no-folders">
      <p style="color: black">No folders in this directory.</p>
    </div>
  </main>
  <main v-else>
    <div class="nav-button-group">
      <p class="current-directory">
        Current shared directory: {{ currentDir }}
      </p>
      <button @click="chooseNewDir" class="nav-button">
        Change shared folder
      </button>
    </div>
  </main>
</template>

<style scoped>
.outer-div {
  margin: auto;
  /* display: flex; */
}

.nav-button-group {
  display: flex;
  margin: auto;
  justify-content: center;
  height: 100%;
  gap: 10px;

  .nav-button {
    background-color: lightgrey;
  }
  .current-directory {
    color: black;
    border: solid black 2px;
    padding: 10px;
    background-color: lightgrey;
    /* margin-top: 0; */
  }
}

.folder-selection-box {
  display: grid;
  grid-template-columns: auto auto auto;
  max-height: 20vh;
  overflow-y: scroll;

  .folder-button {
    width: 100%;
    display: flex;
    flex-direction: row-reverse;
    justify-content: center;
    gap: 5px;

    .folder-image {
      width: 3%;
    }
  }
  /* .no-folders {
    margin: auto;

  } */
}

.no-folders {
  min-height: 30vh;
}

.result {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
}

.input-box .btn {
  width: 60px;
  height: 30px;
  line-height: 30px;
  border-radius: 3px;
  border: none;
  margin: 0 0 0 20px;
  padding: 0 8px;
  cursor: pointer;
}

.input-box .btn:hover {
  background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
  color: #333333;
}

.input-box .input {
  border: none;
  border-radius: 3px;
  outline: none;
  height: 30px;
  line-height: 30px;
  padding: 0 10px;
  background-color: rgba(240, 240, 240, 1);
  -webkit-font-smoothing: antialiased;
}

.input-box .input:hover {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.input-box .input:focus {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}
</style>
