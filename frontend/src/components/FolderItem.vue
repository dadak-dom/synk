<script lang="ts" setup>
import { onMounted, onUpdated, ref } from "vue";

defineProps({
  text: String,
  file: Boolean,
});

const showItem = ref<boolean>(false);

onMounted(() => {
  setTimeout(() => (showItem.value = true), 1);
});

onUpdated(() => {
  showItem.value = false;
  setTimeout(() => (showItem.value = true), 1);
});
</script>

<template>
  <main>
    <button class="folder-button">
      <Transition name="slide-fade">
        <img
          v-show="!file && showItem"
          class="folder-image"
          src="../assets/images/folder.png"
        />
      </Transition>
      <Transition name="slide-fade">
        <div class="folder-text" v-if="showItem">
          {{ text }}
        </div>
      </Transition>
    </button>
  </main>
</template>

<style scoped>
.folder-button {
  width: 100%;
  display: flex;
  justify-content: left;
  gap: 5px;
}

.folder-image {
  width: 10%;
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
</style>
