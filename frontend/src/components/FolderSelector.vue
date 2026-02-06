<script lang="ts" setup>
import { onMounted, ref, useTemplateRef } from "vue";
// import what I need from the backend...
import { FolderSelectorControl } from "../../wailsjs/go/main/App";

const enableSelector = ref(false);
const currentDir = ref();
const foldersInCurrentDir = ref<string[]>([]);
const filesInCurrentDir = ref<string[]>([]);
var rawFoldersInCurrentDir = <string[]>[];

// Options for viewing the file selector:
const showHiddenFiles = ref(false);
const onlyFolders = ref(false); // FIXME : should be named "showFiles", too tired to fix it now :(

// Modal settings
// const modal = useTemplateRef("modal");
const modalOpacity = ref("0");

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
    foldersInCurrentDir.value = value.Folders;
    rawFoldersInCurrentDir = foldersInCurrentDir.value;
    filesInCurrentDir.value = value.Files;
    handleShowHidden();
    handleOnlyFolders();
    console.log(foldersInCurrentDir.value);
  });
}

function moveDownDir(f: string) {
  FolderSelectorControl(
    currentDir.value,
    FolderSelectorCommands.MOVE_DOWN,
    f,
  ).then((value) => {
    currentDir.value = value.Directory;
    foldersInCurrentDir.value = value.Folders;
    rawFoldersInCurrentDir = foldersInCurrentDir.value;
    filesInCurrentDir.value = value.Files;
    handleShowHidden();
    handleOnlyFolders();
    console.log(foldersInCurrentDir.value);
  });
}

function goHome() {
  FolderSelectorControl("", FolderSelectorCommands.GO_HOME, "").then(
    (value) => {
      currentDir.value = value.Directory;
      foldersInCurrentDir.value = value.Folders;
      rawFoldersInCurrentDir = foldersInCurrentDir.value;
      filesInCurrentDir.value = value.Files;
      handleShowHidden();
      handleOnlyFolders();
      console.log(foldersInCurrentDir.value);
    },
  );
}

function selectFolder() {
  modalOpacity.value = "0";
  // enableSelector.value = false;
  setTimeout(() => (enableSelector.value = false), 300);
  FolderSelectorControl(currentDir.value, FolderSelectorCommands.SELECT, "");
}

function cancelFolderSelect() {
  // enableSelector.value = false;
  // if (modal.value !== null) {
  //   modal.value.style.opacity = "0";
  // }
  modalOpacity.value = "0";
  setTimeout(() => (enableSelector.value = false), 300);
}

function handleShowHidden() {
  if (!showHiddenFiles.value) {
    rawFoldersInCurrentDir = foldersInCurrentDir.value;
    foldersInCurrentDir.value = foldersInCurrentDir.value.filter(
      (file) => !file.startsWith("."),
    );
  } else {
    foldersInCurrentDir.value = rawFoldersInCurrentDir;
  }
}

function chooseNewDir() {
  enableSelector.value = true;
  setTimeout(() => (modalOpacity.value = "1"), 300);
  // console.log(modal.value);

  // if (modal.value !== null) {
  //   modal.value.style.opacity = "1";
  // }
  // modalOpacity.value = "1";
}

function handleOnlyFolders() {}

onMounted(() => {
  FolderSelectorControl("", FolderSelectorCommands.INIT, "").then((value) => {
    currentDir.value = value.Directory;
    foldersInCurrentDir.value = value.Folders;
    filesInCurrentDir.value = value.Files;
    handleShowHidden();
    handleOnlyFolders();
    // if (modal.value !== null) {
    //   modal.value.style.opacity = "0"; // init modal opacity to 0
    //   modal.value.style.transition = "opacity 1s";
    // }
    console.log(foldersInCurrentDir.value);
    if (currentDir.value == "") {
      enableSelector.value = true;
    }
  });
});
</script>

<template>
  <main ref="modal" class="folder-selector-modal" v-show="enableSelector">
    <div class="modal-wrapper">
      <div class="options-and-current-dir">
        <p class="current-directory">{{ currentDir }}</p>
        <div class="options-wrapper">
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
          <div class="option">
            <div>Show files</div>
            <label class="switch">
              <input
                type="checkbox"
                v-bind:checked="onlyFolders"
                v-model="onlyFolders"
                @change="handleOnlyFolders"
              />
              <span class="slider round"></span>
            </label>
          </div>
        </div>
      </div>
      <div class="nav-button-group">
        <button class="nav-button" @click="moveUpDir">↑ Move Up</button>
        <button class="nav-button" @click="goHome">Go Home</button>
        <button class="nav-button" @click="cancelFolderSelect">Cancel</button>
        <button class="nav-button" @click="selectFolder">
          Select this folder
        </button>
      </div>

      <div class="folder-selection-box">
        <div
          v-if="
            foldersInCurrentDir == undefined || foldersInCurrentDir.length > 0
          "
          v-for="folder in foldersInCurrentDir"
        >
          <button class="folder-button" @click="moveDownDir(folder)">
            <img class="folder-image" src="../assets/images/folder.png" />
            <div class="folder-text">
              {{ folder }}
            </div>
          </button>
        </div>
        <div v-if="onlyFolders" v-for="file in filesInCurrentDir">
          <button class="folder-button">
            <!-- <img class="folder-image" src="../assets/images/folder.png" /> -->
            <div class="folder-text">
              {{ file }}
            </div>
          </button>
        </div>
      </div>
      <!-- <div v-else class="folder-selection-box">
        <p style="color: black">No folders in this directory.</p>
      </div> -->
    </div>
  </main>
  <main>
    <div class="change-folder-button">
      <img class="" @click="chooseNewDir" src="../assets/images/folder.png" />
    </div>
  </main>
</template>

<style scoped>
/* .outer-div {
  margin: auto;
} */

.folder-selector-modal {
  position: absolute;
  backdrop-filter: blur(10px);
  width: 100%;
  height: 100%;
  opacity: v-bind(modalOpacity);
  transition: opacity 0.3s;
  /* background-color: #142431; */
}

.modal-wrapper {
  position: relative;
  top: 50%;
  left: 30%;
  transform: translate(-25%);
  display: flex;
  justify-content: center;
  height: 40vh;
  width: 80vw;
}

.options-and-current-dir {
  width: 40vw;

  .options-wrapper {
    display: grid;
    grid-template-columns: auto auto auto;
    background-color: lightgrey;
    .option {
      /* background-color: #aeb3ba; */
      color: black;
    }
  }

  .current-directory {
    color: black;
    border: solid black 2px;
    padding: 10px;
    background-color: lightgrey;
  }
}

.nav-button-group {
  display: flex;
  justify-content: center;
  flex-direction: column;
  /* height: 100%; */
  width: 150px;
  gap: 10px;

  .nav-button {
    background-color: lightgrey;
    width: 50%;
    margin: auto;
  }
}

.folder-selection-box {
  display: list-item;
  /* grid-template-columns: auto auto auto; */
  /* max-height: 20vh; */
  /* height: 100%; */
  overflow-y: scroll;
  overflow-x: hidden;
  width: 20vw;

  .folder-button {
    width: 100%;
    display: flex;
    justify-content: left;
    gap: 5px;

    .folder-image {
      width: 10%;
    }
  }
}

.no-folders {
  min-height: 30vh;
  width: 20%;
  background-color: lightgrey;
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
  background-color: #98cb98;
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
