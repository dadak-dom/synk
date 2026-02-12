<script lang="ts" setup>
import { onMounted, onUpdated, ref, useTemplateRef } from "vue";
// import what I need from the backend...
import {
  FolderSelectorControl,
  GetSharedDirectory,
} from "../../wailsjs/go/main/App";
import FolderSelectorItems from "../components/FolderSelectorItems.vue";

const enableSelector = ref(false);
const currentDir = ref();
const foldersInCurrentDir = ref<string[]>([]);
const filesInCurrentDir = ref<string[]>([]);
var rawFoldersInCurrentDir = <string[]>[];
var rawFilesInCurrentDir = <string[]>[];

const sharedDirectory = ref("");

// Options for viewing the file selector:
const showHiddenFiles = ref(false);
const showFiles = ref(false);

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
    rawFilesInCurrentDir = value.Files;
    handleShowHidden();
    handleShowFiles();
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
    rawFilesInCurrentDir = value.Files;
    handleShowHidden();
    handleShowFiles();
    console.log(foldersInCurrentDir.value);
  });
}

function goHome() {
  FolderSelectorControl("", FolderSelectorCommands.GO_HOME, "").then(
    (value) => {
      currentDir.value = value.Directory;
      foldersInCurrentDir.value = value.Folders;
      rawFoldersInCurrentDir = foldersInCurrentDir.value;
      rawFilesInCurrentDir = value.Files;
      filesInCurrentDir.value = value.Files;
      handleShowHidden();
      handleShowFiles();
      console.log(foldersInCurrentDir.value);
    },
  );
}

function selectFolder() {
  openFolderSelector.value = false;
  setTimeout(() => (showFolderButton.value = true), 300);
  FolderSelectorControl(currentDir.value, FolderSelectorCommands.SELECT, "");
  sharedDirectory.value = currentDir.value;
}

function cancelFolderSelectTransition() {
  openFolderSelector.value = false;
  setTimeout(() => {
    showFolderButton.value = true;
  }, 300);
}

function changeFolderTransition() {
  showFolderButton.value = false;
  setTimeout(() => {
    openFolderSelector.value = true;
  }, 300);
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

function handleShowFiles() {
  if (showFiles.value) {
    filesInCurrentDir.value = rawFilesInCurrentDir;
  } else {
    filesInCurrentDir.value = [];
  }
}

function chooseNewDir() {
  enableSelector.value = true;
  console.debug("chooseNewDir event fired...");
  setTimeout(() => (modalOpacity.value = "1"), 300);
}

onMounted(() => {
  FolderSelectorControl("", FolderSelectorCommands.INIT, "").then((value) => {
    console.debug("Running FolderSelectorInit...");
    currentDir.value = value.Directory;
    foldersInCurrentDir.value = value.Folders;
    filesInCurrentDir.value = value.Files;
    rawFilesInCurrentDir = value.Files;
    handleShowHidden();
    handleShowFiles();
    console.log(foldersInCurrentDir.value);
  });
  GetSharedDirectory().then((dir) => {
    sharedDirectory.value = dir;
  });
});

const openFolderSelector = ref<boolean>(false);
const showFolderButton = ref<boolean>(true);
</script>

<template>
  <div class="outer-view">
    <div class="folder-selection-view">
      <div class="folder-view-wrapper">
        <h1>Shared Folder</h1>
        <div class="shared-directory tooltip">
          <span class="tooltiptext">Your shared directory</span
          >{{ sharedDirectory }}
        </div>
        <Transition name="slide-fade">
          <div class="change-folder-button-wrapper" v-if="showFolderButton">
            <button
              class="change-folder-button"
              @click="changeFolderTransition"
              v-if="showFolderButton"
            >
              Change Shared Folder
            </button>
          </div>
        </Transition>
        <Transition name="slide-fade">
          <div class="folder-selector-box" v-if="openFolderSelector">
            <div class="folder-selection-box">
              <FolderSelectorItems
                :folders="foldersInCurrentDir"
                :files="filesInCurrentDir"
                @move-down-dir="moveDownDir"
              />
            </div>
            <div class="nav-button-group">
              <button class="nav-button" @click="cancelFolderSelectTransition">
                Cancel
              </button>
              <button class="nav-button" @click="moveUpDir">↑ Move Up</button>
              <button class="nav-button" @click="goHome">Go Home</button>
              <button class="nav-button" @click="selectFolder">
                Select this folder
              </button>
            </div>
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
                      v-bind:checked="showFiles"
                      v-model="showFiles"
                      @change="handleShowFiles"
                    />
                    <span class="slider round"></span>
                  </label>
                </div>
              </div>
            </div>
          </div>
        </Transition>
      </div>
    </div>
  </div>
</template>

<style scoped>
.folder-selection-view {
  height: 86%;
  width: 100%;
  display: flex;
  justify-content: center;
}

.folder-view-wrapper {
  width: 90%;
}

.folder-selector-box {
  height: 100%;
  width: 100%;
  border-left: 1px solid grey;
  border-right: 1px solid grey;
  background: rgba(20, 20, 20, 0.9);
}

.change-folder-button-wrapper {
  background-color: rgba(30, 30, 30, 0.8);
  border-left: 1px solid grey;
  border-right: 1px solid grey;
  width: 100%;
  height: 90%;
  border-bottom-left-radius: 40px;
  border-bottom-right-radius: 40px;
}

.change-folder-button-wrapper .change-folder-button {
  color: white;
  padding: 20px;
  border: 2px solid rgb(33, 33, 33);
  border-radius: 10px;
  margin-bottom: 40px;
}

.shared-directory {
  /* margin: 10%; */
  background: linear-gradient(
    180deg,
    rgba(148, 148, 148, 0.9) 0,
    rgba(7, 7, 7, 0.93) 20%,
    rgba(19, 19, 19, 0.9) 40%,
    rgba(30, 30, 30, 0.93) 100%
  );
  /* padding: 20px; */
  padding-top: 20px;
  padding-bottom: 20px;
  /* margin-top: 10%; */
  width: 100%;
  border-left: 1px solid grey;
  border-right: 1px solid grey;
  border-top: 1px solid grey;
  border-top-left-radius: 40px;
  border-top-right-radius: 40px;
}

.folder-view-spacer {
  padding-top: 20px;
  background: linear-gradient(
    180deg,
    rgba(7, 7, 7, 0.93) 0,
    rgba(19, 19, 19, 0.9) 20%,
    rgba(71, 71, 71, 0.9) 40%,
    rgba(148, 148, 148, 0.9) 100%
  );
  border-left: 1px solid grey;
  border-right: 1px solid grey;
  width: 100%;
}

.slide-fade-enter-active {
  transition: all 0.3s ease-out;
}

.slide-fade-leave-active {
  transition: all 0.3s ease-in-out;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateX(20px);
  opacity: 0;
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

.options-and-current-dir .options-wrapper {
  display: grid;
  grid-template-columns: auto auto auto;
  /* background-color: lightgrey; */
}

.options-and-current-dir .option {
  color: white;
}

.current-directory {
  color: white;
  /* border: solid black 2px; */
  padding: 10px;
  /* background-color: lightgrey;
   */
}

.nav-button-group {
  display: flex;
  justify-content: center;
  flex-direction: row;
}

.change-folder-button {
  margin-top: 50px;
  background: linear-gradient(
    180deg,
    rgba(148, 148, 148, 0.9) 0,
    rgba(7, 7, 7, 0.93) 20%,
    rgba(19, 19, 19, 0.9) 40%,
    rgba(105, 102, 102, 0.93) 100%
  );
}

.nav-button-group .nav-button {
  background-color: lightgrey;
  /* width: 50%; */
  /* margin: auto; */
  height: 100%;
  color: white;
  background-color: rgba(45, 44, 44, 0.7);
}

.folder-selection-box {
  display: list-item;
  height: 30vh;
  overflow-y: scroll;
  overflow-x: hidden;
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
  /* background-color: #98cb98; */
  /* background-image: url("../assets/images/slider-background.jpg");
   */
  background-color: limegreen;
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

/* Tooltip container */
.tooltip {
  position: relative;
  display: inline-block;
  /* border-bottom: 1px dotted black; Add dots under the hoverable text */
  cursor: pointer;
}

/* Tooltip text */
.tooltiptext {
  visibility: hidden; /* Hidden by default */
  width: 70px;
  background-color: black;
  color: #ffffff;
  text-align: center;
  padding: 5px 0;
  border-radius: 6px;
  position: absolute;
  z-index: 1; /* Ensure tooltip is displayed above content */
  font-size: 10px;
  bottom: 100%;
  left: 35%;
  margin-left: -35%;
}

/* Show the tooltip text on hover */
.tooltip:hover .tooltiptext {
  visibility: visible;
}
</style>
