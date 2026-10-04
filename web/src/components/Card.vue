<template>
  <button
    type="button"
    class="vote-card"
    :class="{ 'vote-card-active': active }"
    :aria-pressed="active"
    :aria-label="`Vote ${value}${active ? ', selected' + confidenceDescription : ''}`"
    :disabled="disabled"
    :aria-disabled="disabled || busy"
    translate="no"
    @click="!busy && $emit('click')"
  >
    <span v-if="active" class="vote-card-confidence" aria-hidden="true">{{ confidence === 'high' ? '!' : confidence === 'low' ? '?' : '' }}</span>
    <span class="vote-card-value">{{ value }}</span>
    <span v-if="active" class="vote-card-indicator" aria-hidden="true">✓</span>
  </button>
</template>

<script>
export default {
  props: {
    active: Boolean,
    value: String,
    confidence: String,
    disabled: Boolean,
    busy: Boolean,
  },

  computed: {
    confidenceDescription() {
      if (this.confidence === 'high') {
        return ', very confident'
      }

      return this.confidence === 'low' ? ', not sure' : ''
    },
  },
}
</script>

<style scoped>
.vote-card {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  position: relative;
  width: var(--vote-card-width, 56px);
  height: var(--vote-card-height, 92px);
  flex-shrink: 0;
  border: 1px solid #bcbcbc;
  border-radius: 12px;
  background: var(--poker-surface);
  color: var(--poker-text);
  font: inherit;
  font-size: 1.2rem;
  font-weight: 600;
  cursor: pointer;
  transform: translateY(0);
  transition:
    transform 240ms cubic-bezier(0.2, 0.8, 0.2, 1),
    box-shadow 240ms ease,
    border-color 180ms ease,
    background 180ms ease,
    opacity 180ms ease var(--pending-fade-delay, 0ms);
}

@media (hover: hover) and (pointer: fine) {
  .vote-card:hover:not(:disabled):not(.vote-card-active) {
    border-color: var(--poker-accent);
    background: var(--poker-accent-soft);
    transform: translateY(-4px);
    box-shadow: 0 3px 6px #00000018;
  }
}

.vote-card-active {
  border: 2px solid var(--poker-accent);
  background: var(--poker-surface);
  transform: translateY(-6px);
  box-shadow: 0 3px 6px #00000018;
}

.vote-card-confidence {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 20px;
  border-radius: 9px 9px 0 0;
  background: var(--poker-accent);
  color: var(--poker-on-accent);
  font-size: 0.8rem;
  line-height: 20px;
  text-align: center;
}

.vote-card-indicator {
  position: absolute;
  right: 6px;
  bottom: 4px;
  color: var(--poker-accent-ink);
  font-size: 0.7rem;
}

.vote-card[aria-disabled="true"] {
  cursor: wait;
}

.vote-card[aria-disabled="true"]:not(:disabled) {
  --pending-fade-delay: 300ms;
  opacity: 0.6;
}

.vote-card:disabled {
  opacity: 0.6;
}

@media (max-width: 600px) {
  .vote-card {
    width: 48px;
    height: 80px;
    font-size: 1.05rem;
  }
}
</style>
