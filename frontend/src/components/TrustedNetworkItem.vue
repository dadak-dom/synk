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
  borderColor: String,
  itemText: String,
});

const emit = defineEmits(
  [
    "removeTrustedNetwork"
  ]
)

function removeTrustedNetwork(network: string) {
  emit("removeTrustedNetwork", network)
}

</script>

<template>
  <div class="network-item-wrapper">
    <div class="network-item-wrapper-name">
      {{itemText}}
    </div>
    <button @click="removeTrustedNetwork(itemText ?? '')">X</button>
  </div>
</template>

<style scoped>
.network-item-wrapper {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 5px;
  border: 1px solid v-bind(borderColor);
  padding: 10px;
  border-radius: 20px;
}
</style>
