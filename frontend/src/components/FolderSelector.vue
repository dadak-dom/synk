<script lang="ts" setup>
import { onMounted, ref } from "vue";
// import what I need from the backend...
import { FolderSelectorControl } from "../../wailsjs/go/main/App";

const enableSelector = ref(false);
const currentDir = ref();
const filesInCurrentDir = ref<string[]>([]);
var rawFilesInCurrentDir = <string[]>[];

// Options for viewing the file selector:
const showHiddenFiles = ref(false);

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
    rawFilesInCurrentDir = filesInCurrentDir.value;
    handleShowHidden();
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
    rawFilesInCurrentDir = filesInCurrentDir.value;
    handleShowHidden();
    console.log(filesInCurrentDir.value);
  });
}

function goHome() {
  FolderSelectorControl("", FolderSelectorCommands.GO_HOME, "").then(
    (value) => {
      currentDir.value = value.Directory;
      filesInCurrentDir.value = value.Files;
      rawFilesInCurrentDir = filesInCurrentDir.value;
      handleShowHidden();
      console.log(filesInCurrentDir.value);
    },
  );
}

function selectFolder() {
  enableSelector.value = false;
  FolderSelectorControl(currentDir.value, FolderSelectorCommands.SELECT, "");
}

function cancelFolderSelect() {
  enableSelector.value = false;
}

function handleShowHidden() {
  if (!showHiddenFiles.value) {
    rawFilesInCurrentDir = filesInCurrentDir.value;
    filesInCurrentDir.value = filesInCurrentDir.value.filter(
      (file) => !file.startsWith("."),
    );
  } else {
    filesInCurrentDir.value = rawFilesInCurrentDir;
  }
}

function chooseNewDir() {
  enableSelector.value = true;
}

onMounted(() => {
  FolderSelectorControl("", FolderSelectorCommands.INIT, "").then((value) => {
    currentDir.value = value.Directory;
    filesInCurrentDir.value = value.Files;
    handleShowHidden();
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
      <button class="nav-button" @click="cancelFolderSelect">Cancel</button>
      <button class="nav-button" @click="selectFolder">
        Select this folder
      </button>
      <div class="option">
        <div>Show hidden folders</div>
        <label class="switch">
          <input
            type="checkbox"
            v-bind:checked="showHiddenFiles"
            v-model="showHiddenFiles"
            @change="handleShowHidden"
          />
          <span class="slider round"></span>
        </label>
      </div>
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
}

.nav-button-group {
  display: flex;
  margin: auto;
  justify-content: center;
  height: 100%;
  gap: 10px;

  .option {
    background-color: #aeb3ba;
    color: black;
  }

  .nav-button {
    background-color: lightgrey;
  }
  .current-directory {
    color: black;
    border: solid black 2px;
    padding: 10px;
    background-color: lightgrey;
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
}

.no-folders {
  min-height: 30vh;
}

.result {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
}

/* The switch - the box around the slider */
.switch {
  position: relative;
  display: inline-block;
  width: 60px;
  height: 34px;
}

/* Hide default HTML checkbox */
.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

/* The slider */
.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  -webkit-transition: 0.4s;
  transition: 0.4s;
}

.slider:before {
  position: absolute;
  content: "";
  height: 26px;
  width: 26px;
  left: 4px;
  bottom: 4px;
  background-color: white;
  -webkit-transition: 0.4s;
  transition: 0.4s;
}

input:checked + .slider {
  background-color: #33f321;
}

input:focus + .slider {
  box-shadow: 0 0 1px #2196f3;
}

input:checked + .slider:before {
  -webkit-transform: translateX(26px);
  -ms-transform: translateX(26px);
  transform: translateX(26px);
}

/* Rounded sliders */
.slider.round {
  border-radius: 34px;
}

.slider.round:before {
  border-radius: 50%;
}
</style>
