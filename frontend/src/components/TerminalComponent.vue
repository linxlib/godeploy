<template>
  <div ref="terminal" class="terminal-container"></div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue';

import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import '@xterm/xterm/css/xterm.css';

const terminal = ref<HTMLDivElement | null>(null);
const xterm = new Terminal();
const fitAddon = new FitAddon();

onMounted(() => {
  if (terminal.value) {
    xterm.loadAddon(fitAddon);
    xterm.open(terminal.value);
    fitAddon.fit();

    // Example: Writing to the terminal
    xterm.write('\x1b[31mWelcome to the xterm.js terminal!\x1b[0m\r\n');

    // Handle window resize
    window.addEventListener('resize', () => fitAddon.fit());
  }
})

defineExpose({
  xterm
})


</script>

<style scoped>
.terminal-container {
  width: 100%;
  height: 100%;
}
</style>