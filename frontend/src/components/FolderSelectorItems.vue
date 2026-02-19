<script lang="ts" setup>
// TODO: See if I can abstract this out (New Feature 1)
// then I can use the same component for the "synk-ignore" feature, maybe by passing in a custom function on click
import FolderItem from "./FolderItem.vue";

defineProps({
  folders: { type: Array<string>, required: true },
  files: { type: Array<string>, required: true },
  folderFunc: { type: Function, required: true },
  fileFunc: { type: Function, required: true },
  ignoredFiles: { type: Array<string>, required: true },
  ignoreFolders: { type: Array<string>, required: true },
});

const emit = defineEmits(["moveDownDir"]);

function folderClick(folder: string) {
  emit("moveDownDir", folder);
}
</script>

<template>
  <main>
    <div v-if="folders.length > 0" v-for="folder in folders">
      <!-- <FolderItem :text="folder" @click="folderClick(folder)" :file="false" /> -->
      <FolderItem
        :text="folder"
        @click="folderFunc(folder)"
        :file="false"
        :ignored="ignoreFolders.indexOf(folder) != -1"
      />
    </div>
    <div v-if="files != undefined" v-for="file in files">
      <FolderItem
        :text="file"
        @click="fileFunc(file)"
        :file="true"
        :ignored="ignoredFiles.indexOf(file) != -1"
      />
    </div>
  </main>
</template>

<style scoped></style>
