<script lang="ts" setup>
import { onMounted, onUpdated, ref } from "vue";

defineProps({
  text: String,
  file: Boolean,
  ignored: Boolean,
  firstColor: String,
  secondColor: String,
  thirdColor: String,
  fourthColor: String,
  textColor: String,
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
      <Transition name="fade">
        <div class="folder-text" v-if="showItem">
          <div v-if="ignored" style="color: red">
            <s>{{ text }}</s>
          </div>
          <div v-else>{{ text }}</div>
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
  /* background: linear-gradient(
    300deg,
    rgba(86, 86, 86, 0.9) 0,
    rgba(30, 30, 30, 0.93) 1%
  ); */
  background: linear-gradient(
    300deg,
    v-bind(firstColor) 0,
    v-bind(secondColor) 5%,
    v-bind(thirdColor) 20%,
    v-bind(fourthColor) 100%
  );
  border: none;
  border-radius: 1px;
  /* color: white;
   */
  color: v-bind(textColor);
  cursor: pointer;
}

.folder-image {
  width: 5%;
}
.fade-enter-active {
  transition: all 0.3s ease-out;
}

.fade-leave-active {
  transition: all 0.3s ease-in-out;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
