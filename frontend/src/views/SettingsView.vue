<script lang="ts" setup>
import { inject, onMounted, ref } from "vue";
import { GetLocalIP } from "../../wailsjs/go/main/App";
import { Theme, AllThemeNames } from "../interfaces/theme";
const myPrivateIP = ref<string>("");
const userTheme = inject<Theme>("theme");
const emit = defineEmits(["updateTheme"]);
const selectedTheme = ref<string>(userTheme?.name ? userTheme.name : "");

function themeInput(th: string) {
  emit("updateTheme", th);
  console.log("change event", th);
}

onMounted(() => {
  GetLocalIP().then((ip) => {
    myPrivateIP.value = ip;
  });
});
</script>

<template>
  <main>
    <div class="title">
      <h1>Settings</h1>
    </div>
    <div class="settings">
      <p>My private IP: {{ myPrivateIP }}</p>
      <div class="theme-selector">
        <label for="themes">Theme:</label>
        <select
          name="themes"
          v-model="selectedTheme"
          @change="themeInput(selectedTheme)"
        >
          <template v-for="theme in AllThemeNames">
            <option :value="theme">{{ theme }}</option>
            <p>{{ userTheme?.name }}</p>
          </template>
        </select>
      </div>
    </div>
  </main>
</template>

<style scoped>
.settings {
  margin-top: 20px;
}
.theme-selector {
  display: flex;
  justify-content: center;
  gap: 10px;
}
</style>
