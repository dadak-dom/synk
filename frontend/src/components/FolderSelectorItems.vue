<script lang="ts" setup>
import FolderItem from "./FolderItem.vue";

defineProps({
  folders: Array<string>,
  files: Array<string>,
});

const emit = defineEmits(["moveDownDir"]);

function folderClick(folder: string) {
  emit("moveDownDir", folder);
}
</script>

<template>
  <main>
    <!-- <div v-for="folder in folders">test: {{ folder }}</div> -->
    <div
      v-if="folders == undefined || folders.length > 0"
      v-for="folder in folders"
    >
      <FolderItem :text="folder" @click="folderClick(folder)" :file="false" />
    </div>
    <div v-if="files != undefined" v-for="file in files">
      <FolderItem :text="file" :file="true" />
    </div>
  </main>
</template>

<style scoped>
.folder-button {
  width: 100%;
  display: flex;
  justify-content: left;
  gap: 5px;

  .folder-image {
    width: 10%;
  }
}
</style>
