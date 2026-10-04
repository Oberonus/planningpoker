<template>
  <div class="home-page">
    <header class="home-header">
      <router-link class="brand" to="/">
        <img class="brand-mark" :src="$router.options.base + 'logo.svg'" alt="" aria-hidden="true" width="36" height="36">
        Planning Poker
      </router-link>
    </header>

    <main class="home-content">
      <section class="home-intro" aria-labelledby="home-title">
        <h1 id="home-title">{{ welcomeName ? `Hi ${welcomeName}!` : 'Let’s estimate.' }}</h1>
        <p class="home-description">Create a game and send the link to your team. Everyone picks a card, then you reveal them together.</p>

        <div class="sample-deck" aria-hidden="true">
          <span class="sample-card sample-card-left">3</span>
          <span class="sample-card sample-card-center">5</span>
          <span class="sample-card sample-card-right">8</span>
        </div>

      </section>

      <section class="setup-panel" aria-labelledby="setup-title">
        <h2 id="setup-title">Start a game</h2>
        <p class="muted setup-description">Choose your name and a card deck.</p>

        <v-form ref="setupForm" @submit.prevent="startGame">
          <v-text-field
            v-model="displayName"
            label="Your name"
            autocomplete="nickname"
            :rules="nameRules"
            :disabled="creating"
            outlined
            dense
          ></v-text-field>

          <v-text-field
            v-model="gameName"
            label="Game name (optional)"
            placeholder="e.g. Sprint planning"
            :disabled="creating"
            outlined
            dense
          ></v-text-field>

          <v-text-field
            v-model="ticketURL"
            label="Ticket URL (optional)"
            placeholder="https://"
            :rules="urlRules"
            :disabled="creating"
            outlined
            dense
          ></v-text-field>

          <v-select
            v-model="deck"
            :items="decks"
            item-text="name"
            label="Card deck"
            return-object
            :disabled="creating"
            outlined
            dense
            hide-details
          ></v-select>
          <p class="deck-preview muted" aria-live="polite">{{ deck.types.join(' · ') }}</p>

          <v-alert v-if="error" type="error" text dense role="alert">{{ error }}</v-alert>

          <v-btn type="submit" color="primary" large block depressed :loading="creating">
            Create game
            <v-icon right small aria-hidden="true">mdi-arrow-right</v-icon>
          </v-btn>
        </v-form>
        <p class="setup-footnote muted">Everyone can vote and reveal the cards.</p>
      </section>
    </main>
  </div>
</template>

<script>
import game from "@/models/game"
import user from "@/models/user"
import { ticketURLRules } from "@/models/validation"

export default {
  name: 'Home',

  data: () => ({
    welcomeName: user.name || '',
    displayName: user.name || '',
    gameName: '',
    ticketURL: '',
    deck: game.decks[0],
    decks: game.decks,
    creating: false,
    error: '',
    nameRules: [value => !!value.trim() || 'Enter your name to join the game.'],
    urlRules: ticketURLRules,
  }),

  methods: {
    async startGame() {
      if (this.creating || !this.$refs.setupForm.validate()) {
        return
      }

      this.creating = true
      this.error = ''

      try {
        const name = this.displayName.trim()
        user.name = name
        await user.authenticate()

        if (user.name !== name) {
          await user.update(name)
        }

        const gameID = await game.create(this.gameName.trim(), this.ticketURL.trim(), this.deck)
        await this.$router.push({name: 'Games', params: {id: gameID}})
      } catch (error) {
        this.error = "We couldn't create your game. Check your connection and try again."
      } finally {
        this.creating = false
      }
    },
  },
}
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 48px;
}

.home-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
}

.home-content {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(0, 1fr);
  align-items: center;
  gap: 72px;
  padding: 72px 0 40px;
}

.home-intro h1 {
  font-size: clamp(2.2rem, 4vw, 3.3rem);
  font-weight: 700;
  letter-spacing: -0.025em;
  line-height: 1.15;
  overflow-wrap: anywhere;
}

.home-description {
  margin-top: 24px;
  color: var(--poker-muted);
  font-size: 1.05rem;
  line-height: 1.8;
}

.sample-deck {
  display: flex;
  align-items: center;
  justify-content: center;
  max-width: 320px;
  height: 210px;
  margin: 24px 0;
}

.sample-card {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 86px;
  height: 128px;
  background: var(--poker-surface);
  border: 1px solid var(--poker-border);
  border-radius: 14px;
  color: var(--poker-muted);
  font-size: 2rem;
  font-weight: 600;
}

.sample-card-left {
  transform: rotate(-12deg) translate(8px, 12px);
}

.sample-card-center {
  z-index: 1;
  background: var(--poker-accent);
  border-color: var(--poker-accent);
  color: var(--poker-on-accent);
  transform: translateY(-8px);
  box-shadow: 0 12px 32px #0000000d;
}

.sample-card-right {
  transform: rotate(12deg) translate(-8px, 12px);
}

.setup-panel {
  min-width: 0;
  padding: 32px;
  background: var(--poker-surface);
  border: 1px solid var(--poker-border);
  border-radius: var(--poker-radius);
  box-shadow: 0 8px 40px #00000005;
}

.setup-panel h2 {
  font-size: 1.5rem;
  letter-spacing: -0.025em;
}

.setup-description {
  margin: 8px 0 28px;
  font-size: 0.9rem;
}

.deck-preview {
  /* Reserve two lines so changing decks doesn't move the form actions. */
  min-height: 3.2em;
  margin: 12px 0 20px;
  font-size: 0.8rem;
  line-height: 1.6;
}

.setup-footnote {
  margin: 16px 0 0;
  text-align: center;
  font-size: 0.75rem;
}

@media (max-width: 800px) {
  .home-page {
    padding: 24px;
  }

  .home-content {
    grid-template-columns: minmax(0, 1fr);
    gap: 32px;
    max-width: 480px;
    margin: 0 auto;
    padding-top: 48px;
  }

  .sample-deck {
    display: none;
  }

  .home-description {
    margin-bottom: 0;
  }
}

@media (max-width: 400px) {
  .home-page {
    padding: 20px 16px;
  }

  .setup-panel {
    padding: 24px 20px;
  }
}
</style>
