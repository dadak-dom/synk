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

<style scoped></style>
