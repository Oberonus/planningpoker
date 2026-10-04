<template>
  <v-app>
    <router-view :key="$route.path"></router-view>
  </v-app>
</template>

<script>
export default {
  name: 'App',
};
</script>

<style>
:root {
  --poker-background: #f0f0f0;
  --poker-surface: #ffffff;
  --poker-text: #242424;
  --poker-muted: #646464;
  --poker-border: #dedede;
  --poker-accent: #16a382;
  --poker-on-accent: #ffffff;
  --poker-accent-ink: #126b52;
  --poker-accent-soft: #e8faf5;
  --poker-radius: 20px;
}

html {
  overflow-y: auto;
}

body {
  margin: 0;
  background: var(--poker-background);
}

.theme--light.v-application {
  background: var(--poker-background);
  color: var(--poker-text);
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
}

.v-application .v-btn {
  border-radius: 10px;
  font-weight: 600;
  letter-spacing: 0;
  text-transform: none;
}

.v-application .v-btn.primary {
  color: var(--poker-on-accent) !important;
}

.v-application .v-btn.primary--text {
  color: var(--poker-accent-ink) !important;
  caret-color: var(--poker-accent-ink) !important;
}

.v-application .v-input .primary--text {
  color: var(--poker-accent-ink) !important;
  caret-color: var(--poker-accent-ink) !important;
}

.v-application .v-btn:focus-visible,
.v-application a:focus-visible,
.v-application button:focus-visible {
  outline: 3px solid var(--poker-accent-ink);
  outline-offset: 4px;
}

.v-application .v-dialog {
  border-radius: var(--poker-radius);
}

.dialog-card {
  padding: 12px;
}

.v-application .dialog-card .v-card__title {
  color: var(--poker-text);
  font-size: 1.35rem;
  font-weight: 700;
}

.dialog-description {
  margin: 0 0 24px;
  color: var(--poker-muted);
  line-height: 1.6;
}

.dialog-card .v-card__actions {
  gap: 8px;
  padding: 8px 24px 16px;
}

.brand {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  color: var(--poker-text) !important;
  font-size: 1rem;
  font-weight: 700;
  text-decoration: none;
  white-space: nowrap;
}

.brand-mark {
  display: block;
  flex-shrink: 0;
  width: 36px;
  height: 36px;
}

.eyebrow {
  margin-bottom: 12px;
  color: var(--poker-accent-ink);
  font-size: 0.75rem;
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.muted {
  color: var(--poker-muted);
}

@media (max-width: 600px) {
  .v-application .v-snack__wrapper {
    min-width: 0;
    max-width: calc(100vw - 32px);
  }
}

</style>
