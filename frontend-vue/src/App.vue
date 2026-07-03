<template>
  <router-view />
</template>

<script setup>
</script>

<style>
@import './styles/design-tokens.css';

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: var(--font-family-sans);
  background-color: var(--color-bg);
  color: var(--color-text-primary);
  min-height: 100vh;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  min-height: 100vh;
}

::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: var(--color-bg);
}

::-webkit-scrollbar-thumb {
  background: var(--color-border);
  border-radius: var(--radius-sm);
}

::-webkit-scrollbar-thumb:hover {
  background: var(--color-text-muted);
}

.el-card {
  border: none;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
}

.el-card__header {
  padding: var(--spacing-6);
  border-bottom: 1px solid var(--color-border-light);
}

.el-card__body {
  padding: var(--spacing-6);
}

.el-btn {
  border-radius: var(--radius-md);
  font-weight: var(--font-weight-medium);
}

.el-btn--primary {
  background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-light) 100%);
  border: none;
}

.el-btn--primary:hover {
  background: linear-gradient(135deg, var(--color-accent-dark) 0%, var(--color-accent) 100%);
}

.el-input__wrapper {
  border-radius: var(--radius-md);
  box-shadow: none;
}

.el-input__wrapper:hover {
  box-shadow: var(--shadow-sm);
}

.el-menu {
  border-right: none;
}

.el-menu-item {
  margin: var(--spacing-2) var(--spacing-3);
  border-radius: var(--radius-md);
  transition: all var(--transition-normal);
}

.el-menu-item:hover {
  background-color: rgba(255, 255, 255, 0.08);
}

.el-menu-item.is-active {
  background-color: var(--color-accent);
  color: white;
}

.el-aside {
  position: fixed;
  left: 0;
  top: 0;
  height: 100vh;
  z-index: 100;
}

.el-container {
  margin-left: var(--sidebar-width);
}</style>