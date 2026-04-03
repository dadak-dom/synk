<script lang="ts" setup>
import { inject, onMounted, ref } from "vue";
// import what I need from the backend...
// FIXME: Make it so that I can't select a new folder unless I've moved out of the original dir
import {
  FolderSelectorControl,
  GetConfigValueString,
  GetConfigValueStringList,
  GetSharedDirectory,
  SetConfigItemStringList,
  SetConfigItemString,
} from "../../wailsjs/go/main/App";
import FolderSelectorItems from "../components/FolderSelectorItems.vue";

// Folder and file data
const enableSelector = ref(false);
const currentDir = ref();
const foldersInCurrentDir = ref<string[]>([]);
const filesInCurrentDir = ref<string[]>([]);
const rawFoldersInCurrentDir = ref<string[]>([]); // FIXME: I believe "raw" just mean that it includes hidden folders. Verify this.
const rawFilesInCurrentDir = ref<string[]>([]);

const sharedDirectory = ref("");

// Ignore list
const fileIgnoreList = ref<string[]>([]);
const folderIgnoreList = ref<string[]>([]);
const autoIgnoreAll = ref<boolean>(inject("autoIgnoreAll") ?? false);
// const autoIgnoreAll = ref<boolean>(false);

// Options for viewing the file selector:
const showHiddenFiles = ref(false);
const showFiles = ref(false);

const FolderSelectorCommands = {
  MOVE_UP: 0,
  MOVE_DOWN: 1,
  GO_HOME: 2,
  INIT: 3,
  SELECT: 4,
  CANCEL: 5,
};

const ConfigItems: { [x: string]: string } = {
  SHARED_DIR: "shared_directory.txt",
  FILE_IGNORE: "file_ignore.jsonl",
  FOLDER_IGNORE: "folder_ignore.jsonl",
};

function moveUpDir() {
  FolderSelectorControl(
    currentDir.value,
    FolderSelectorCommands.MOVE_UP,
    "",
  ).then((value) => {
    currentDir.value = value.Directory;
    foldersInCurrentDir.value = value.Folders;
    rawFoldersInCurrentDir.value = foldersInCurrentDir.value;
    filesInCurrentDir.value = value.Files;
    rawFilesInCurrentDir.value = value.Files;
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
    rawFoldersInCurrentDir.value = foldersInCurrentDir.value;
    filesInCurrentDir.value = value.Files;
    rawFilesInCurrentDir.value = value.Files;
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
      rawFoldersInCurrentDir.value = foldersInCurrentDir.value;
      rawFilesInCurrentDir.value = value.Files;
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
  // reset the ignored files list
  SetConfigItemStringList(ConfigItems.FILE_IGNORE, []);
  SetConfigItemStringList(ConfigItems.FOLDER_IGNORE, []);
  folderIgnoreList.value = [];
  fileIgnoreList.value = [];
}

function cancelFolderSelectTransition() {
  // add something to reset the ignore list
  FolderSelectorControl("", FolderSelectorCommands.CANCEL, "").then((value) => {
    rawFilesInCurrentDir.value = value.Files;
    rawFoldersInCurrentDir.value = value.Folders;
  });
  handleShowFiles();
  handleShowHidden();
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
    rawFoldersInCurrentDir.value = foldersInCurrentDir.value;
    foldersInCurrentDir.value = foldersInCurrentDir.value.filter(
      (file) => !file.startsWith("."),
    );
  } else {
    foldersInCurrentDir.value = rawFoldersInCurrentDir.value;
  }
}

function handleShowFiles() {
  if (showFiles.value) {
    filesInCurrentDir.value = rawFilesInCurrentDir.value;
  } else {
    filesInCurrentDir.value = [];
  }
}

function updateFileIgnoreList(file: string) {
  // if file not in the list, add it
  // otherwise, remove it
  console.log("file ignore list: ", fileIgnoreList.value);
  let i = fileIgnoreList.value.indexOf(file);
  if (i == -1) {
    fileIgnoreList.value.push(file);
  } else {
    fileIgnoreList.value = fileIgnoreList.value.filter(
      (value, index) => value != file,
    );
  }
  // remove empty string
  fileIgnoreList.value = fileIgnoreList.value.filter((v) => v != "");
  SetConfigItemStringList(ConfigItems.FILE_IGNORE, fileIgnoreList.value);
  console.log(fileIgnoreList.value);
  // send config to remote peers
  if (peers !== undefined) {
    const peerlist = peers?._rawValue;
    peerlist.forEach((p: string) => {
      const url = "http://" + p + ":8080/updateFileIgnoreList";
      fetch(url, {
        method: "POST",
        body: JSON.stringify(fileIgnoreList.value),
        headers: { "Content-Type": "application/json" },
      });
    });
  }
}

function updateFolderIgnoreList(folder: string) {
  // if file not in the list, add it
  // otherwise, remove it
  console.log("folder ignore list: ", folderIgnoreList.value);
  let i = folderIgnoreList.value.indexOf(folder);
  if (i == -1) {
    folderIgnoreList.value.push(folder);
  } else {
    folderIgnoreList.value = folderIgnoreList.value.filter(
      (value, index) => value != folder,
    );
  }
  // remove empty string
  folderIgnoreList.value = folderIgnoreList.value.filter((v) => v != "");
  SetConfigItemStringList(ConfigItems.FOLDER_IGNORE, folderIgnoreList.value);
  // send config to remote peers
  if (peers !== undefined) {
    const peerlist = peers?._rawValue;
    peerlist.forEach((p: string) => {
      const url = "http://" + p + ":8080/updateFolderIgnoreList";
      fetch(url, {
        method: "POST",
        body: JSON.stringify(folderIgnoreList.value),
        headers: { "Content-Type": "application/json" },
      });
    });
  }
  console.log(folderIgnoreList.value);
}

onMounted(() => {
  FolderSelectorControl("", FolderSelectorCommands.INIT, "").then((value) => {
    console.debug("Running FolderSelectorInit...");
    currentDir.value = value.Directory;
    foldersInCurrentDir.value = value.Folders;
    filesInCurrentDir.value = value.Files;
    rawFilesInCurrentDir.value = value.Files;
    rawFoldersInCurrentDir.value = value.Folders;
    handleShowHidden();
    handleShowFiles();
    console.log(foldersInCurrentDir.value);
    console.log("raw files and folders: ", rawFilesInCurrentDir);
  });
  // FIXME: can be cleaned up using generic functions
  GetSharedDirectory().then((dir) => {
    sharedDirectory.value = dir;
  });
  // Get previously saved ignore lists
  GetConfigValueStringList("file_ignore.jsonl").then((value) => {
    fileIgnoreList.value = value;
  });
  GetConfigValueStringList("folder_ignore.jsonl").then((value) => {
    console.log("Running get folder ignore");
    folderIgnoreList.value = value;
    console.log(folderIgnoreList.value);
  });
});

const openFolderSelector = ref<boolean>(false);
const showFolderButton = ref<boolean>(true);

const peers = inject<any>("peers"); // really should be type: stirng[], but for some reason that breaks when using ._rawValue

// wipe the ignores
function resetIgnores() {
  folderIgnoreList.value = [];
  fileIgnoreList.value = [];
  SetConfigItemStringList(ConfigItems.FOLDER_IGNORE, folderIgnoreList.value);
  SetConfigItemStringList(ConfigItems.FILE_IGNORE, fileIgnoreList.value);
  console.log("peers: ", typeof peers, peers);
  console.log(peers._rawValue);
  const peerlist = peers?._rawValue;
  if (peers !== undefined && peers !== null) {
    peerlist.forEach((p: string) => {
      const url = "http://" + p + ":8080/resetIgnoreList";
      console.log("url: ", url);
      fetch(url);
    });
  }
}

// ignore everything in shared directory
function ignoreAll() {
  // reset everything first, then ignore everything
  console.log(
    "TESTING!!!!: ",
    foldersInCurrentDir.value.length,
    filesInCurrentDir.value.length,
  );
  resetIgnores();
  rawFilesInCurrentDir.value.forEach((value, index) =>
    updateFileIgnoreList(value),
  );
  rawFoldersInCurrentDir.value.forEach((value, index) =>
    updateFolderIgnoreList(value),
  );
}

// toggles the value of the autoIgnoreAll config item
function toggleAutoIgnoreAll() {
  console.log("Auto ignore all: ", autoIgnoreAll.value, !autoIgnoreAll.value);
  SetConfigItemString("auto_ignore_all.txt", autoIgnoreAll.value.toString());
  // if (autoIgnoreAll.value) {
  //   SetConfigItemStringList(ConfigItems.FILE_IGNORE, ["*"]);
  //   SetConfigItemStringList(ConfigItems.FOLDER_IGNORE, ["*"]); // FIXME: If "*" is the folder/file ignore list, have everything redded out in the list
  // } else {
  //   SetConfigItemStringList(ConfigItems.FILE_IGNORE, []);
  //   SetConfigItemStringList(ConfigItems.FOLDER_IGNORE, []); // FIXME: just set to empty for now
  // }
}

// Theme piping
const props = defineProps([
  "firstColor",
  "secondColor",
  "thirdColor",
  "fourthColor",
  "textColor",
  "backgroundColor",
  "borderColor",
]);
</script>

<template>
  <div class="outer-view">
    <div class="folder-selection-view">
      <div class="folder-view-wrapper">
        <h1 class="title">Shared Folder</h1>
        <div class="shared-directory tooltip" v-if="showFolderButton">
          <span class="tooltiptext">Your shared directory</span>
          <div v-if="sharedDirectory != ''">{{ sharedDirectory }}</div>
          <div v-else style="color: red">No shared directory set.</div>
        </div>
        <div class="current-directory shared-directory" v-else>
          {{ currentDir }}
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
            <div class="ignore-list">
              <div>Click on file or folder to have Synk ignore it.</div>
              <FolderSelectorItems
                :folders="rawFoldersInCurrentDir"
                :files="rawFilesInCurrentDir"
                :folder-func="
                  (folder: string) => {
                    updateFolderIgnoreList(folder);
                  }
                "
                :ignore-folders="folderIgnoreList"
                :ignored-files="fileIgnoreList"
                :file-func="
                  (file: string) => {
                    updateFileIgnoreList(file);
                  }
                "
                :receive-only="autoIgnoreAll"
              />
            </div>
            <button class="change-folder-button" @click="resetIgnores">
              Reset Ignores
            </button>
            <button
              class="change-folder-button"
              @click="ignoreAll"
              :disabled="
                folderIgnoreList.length == rawFoldersInCurrentDir.length &&
                fileIgnoreList.length == rawFilesInCurrentDir.length
              "
            >
              Ignore All
            </button>

            <div class="option">
              <div>Receive-Only Mode</div>
              <label class="switch">
                <input
                  type="checkbox"
                  v-bind:checked="autoIgnoreAll"
                  v-model="autoIgnoreAll"
                  @change="toggleAutoIgnoreAll"
                />
                <span class="slider round"></span>
              </label>
            </div>
          </div>
        </Transition>
        <Transition name="slide-fade">
          <div class="folder-selector-box" v-if="openFolderSelector">
            <div class="folder-selection-box">
              <FolderSelectorItems
                :folders="foldersInCurrentDir"
                :files="filesInCurrentDir"
                :ignore-folders="[]"
                :ignored-files="[]"
                :folder-func="
                  (folder: string) => {
                    moveDownDir(folder);
                  }
                "
                :file-func="(file: string) => {}"
                :firstColor="firstColor"
                :secondColor="secondColor"
                :thirdColor="thirdColor"
                :fourthColor="fourthColor"
                :textColor="textColor"
                :receive-only="autoIgnoreAll"
                :borderColor="borderColor"
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

.ignore-list {
  overflow-y: scroll;
  overflow-x: hidden;
  max-height: 200px;
}

.folder-view-wrapper {
  color: v-bind(textColor);
  width: 90%;
}

.outer-view {
  background: v-bind(backgroundColor);
}

.title {
  border-radius: 20px;
  border: 1px solid v-bind(borderColor);
}

.folder-selector-box {
  height: 100%;
  width: 100%;
  border-left: 1px solid v-bind(borderColor);
  border-right: 1px solid v-bind(borderColor);
  background: v-bind(backgroundColor);
}

.change-folder-button-wrapper {
  /* background-color: rgba(30, 30, 30, 0.8); */
  /* background-color: linear-gradient(
    200deg,
    v-bind(firstColor),
    v-bind(fourthColor)
  ); */
  background: v-bind(backgroundColor);
  border-left: 1px solid v-bind(borderColor);
  border-right: 1px solid v-bind(borderColor);
  width: 100%;
  height: 90%;
  border-bottom-left-radius: 40px;
  border-bottom-right-radius: 40px;
}

.change-folder-button-wrapper .change-folder-button {
  padding: 20px;
  border: 2px solid rgb(33, 33, 33);
  border-radius: 10px;
  margin-bottom: 40px;
  margin-top: 10px;
  background: v-bind(backgroundColor);
  color: v-bind(textColor);
}

.shared-directory {
  /* background: linear-gradient(
    200deg,
    v-bind(firstColor) 0,
    v-bind(secondColor) 20%,
    v-bind(thirdColor) 40%,
    v-bind(fourthColor) 100%
  ); */
  background: v-bind(backgroundColor);
  padding-top: 20px;
  padding-bottom: 20px;
  width: 100%;
  border-left: 1px solid v-bind(borderColor);
  border-right: 1px solid v-bind(borderColor);
  border-top: 1px solid v-bind(borderColor);
  border-top-left-radius: 40px;
  border-top-right-radius: 40px;
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
  background: v-bind(backgroundColor);
}

.options-and-current-dir .option {
  /* color: white; */
}

.current-directory {
  animation: blink 0.5s ease-in-out;
  animation-iteration-count: infinite;
  border-width: 1px;
  border-style: solid;
}

@keyframes blink {
  0% {
    border-color: black;
  }
  50% {
    border-color: white;
  }
  100% {
    border-color: black;
  }
}

.nav-button-group {
  display: flex;
  justify-content: center;
  flex-direction: row;
}

.nav-button-group .nav-button {
  height: 100%;
  background: linear-gradient(v-bind(firstColor), v-bind(thirdColor));
  color: v-bind(textColor);
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
