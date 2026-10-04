<template>
  <div class="player-card">
    <div class="player-card-body" :class="cardClass" :aria-label="cardDescription" role="img">
      <span v-if="!card" class="waiting-dot" aria-hidden="true"></span>
      <template v-else-if="card !== '*'">
        <span class="confidence-mark" :title="status" aria-hidden="true">{{ confidence === 'high' ? '!' : confidence === 'low' ? '?' : '' }}</span>
        <span class="revealed-value" translate="no" aria-hidden="true">{{ card }}</span>
      </template>
    </div>
    <span class="player-name" :title="name">{{ name }}</span>
    <span class="player-status" :class="{ 'player-status-voted': !!card }">{{ status }}</span>
  </div>
</template>

<script>
export default {
  props: {
    card: String,
    name: String,
    confidence: String,
    finished: Boolean,
  },

  computed: {
    cardClass() {
      return {
        'player-card-waiting': !this.card,
        'player-card-hidden': this.card === '*',
        'player-card-revealed': !!this.card && this.card !== '*',
      }
    },

    status() {
      if (!this.card) {
        return this.finished ? 'No vote' : 'Thinking…'
      }

      if (this.card === '*') {
        return 'Voted'
      }

      if (this.confidence === 'high') {
        return 'Very confident'
      }

      if (this.confidence === 'low') {
        return 'Not sure'
      }

      return 'Revealed'
    },

    cardDescription() {
      if (!this.card) {
        return `${this.name}: ${this.finished ? 'did not vote' : 'waiting for a vote'}`
      }

      return this.card === '*'
        ? `${this.name}: voted, estimate hidden`
        : `${this.name}: ${this.card}, ${this.status.toLowerCase()}`
    },
  },
}
</script>

<style scoped>
.player-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 104px;
  gap: 6px;
}

.player-card-body {
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  width: 68px;
  height: 112px;
  margin-bottom: 8px;
  border-radius: 13px;
}

.player-card-waiting {
  border: 1px dashed #bcbcbc;
  background: #e8e9ea;
}

.waiting-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #bcbcbc;
}

.player-card-hidden {
  border: 1px solid var(--poker-accent);
  background:
    linear-gradient(-135deg, rgb(35, 210, 170) 10%, transparent),
    repeating-linear-gradient(45deg, rgba(35, 210, 170, 1) 0%, rgba(35, 150, 100, 0.6) 5%, transparent 5%, transparent 10%),
    repeating-linear-gradient(-45deg, rgba(35, 210, 170, 0.4) 0%, rgba(35, 150, 100, 0.5) 5%, transparent 5%, transparent 10%);
  background-color: rgba(35, 210, 170, 0.25);
}

.player-card-revealed {
  border: 2px solid var(--poker-accent);
  background: var(--poker-surface);
  color: var(--poker-text);
}

.confidence-mark {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 24px;
  border-radius: 10px 10px 0 0;
  background: var(--poker-accent);
  color: var(--poker-on-accent);
  font-size: 1.05rem;
  font-weight: 700;
  line-height: 24px;
  text-align: center;
}

.revealed-value {
  font-size: 1.6rem;
  font-weight: 700;
}

.player-name {
  max-width: 100%;
  color: var(--poker-text);
  font-size: 0.875rem;
  font-weight: 600;
  line-height: 1.4;
  text-align: center;
  overflow-wrap: anywhere;
}

.player-status {
  color: var(--poker-muted);
  font-size: 0.7rem;
  text-align: center;
}

.player-status-voted {
  color: var(--poker-accent-ink);
}

@media (max-width: 600px) {
  .player-card {
    width: 88px;
  }
}
</style>
