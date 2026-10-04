<template>
  <div class="game-page">
    <header class="game-header">
      <router-link class="brand" to="/" aria-label="Planning Poker home">
        <img class="brand-mark" :src="$router.options.base + 'logo.svg'" alt="" aria-hidden="true" width="36" height="36">
        <span class="game-brand-name">Planning Poker</span>
      </router-link>

      <div v-if="ready" class="session-title">
        <h1 :title="state.name || state.ticket_url || 'Planning session'">
          <a v-if="ticketLink" :href="ticketLink" target="_blank" rel="noopener noreferrer">
            {{ state.name || state.ticket_url }}
            <v-icon small aria-label="Opens in a new tab">mdi-open-in-new</v-icon>
          </a>
          <span v-else>{{ state.name || 'Planning session' }}</span>
        </h1>
        <span class="session-deck muted">{{ state.cards_deck.name }}</span>
      </div>

      <div class="header-actions">
        <v-btn v-if="ready" color="primary" outlined :disabled="!connected" @click="copyLink">
          <v-icon left small aria-hidden="true">mdi-link-variant</v-icon>
          Invite<span class="invite-label"> players</span>
        </v-btn>
        <v-btn v-if="ready" icon aria-label="Game settings" title="Game settings" :disabled="!connected" @click="changeGameParams">
          <v-icon>mdi-cog-outline</v-icon>
        </v-btn>
        <v-btn v-if="ready" icon :aria-label="`Change your name (${user.name})`" title="Change your name" :disabled="!connected" @click="changeName">
          <v-icon>mdi-account-outline</v-icon>
        </v-btn>
      </div>
    </header>

    <main class="game-content">
      <v-alert v-if="notifier.status === notifier.STATUS_RECONNECTING" type="warning" text role="status">
        Connection lost. Reconnecting… Your game will update when you're back online.
      </v-alert>
      <v-alert v-if="joinFailed" type="error" text role="alert">
        We couldn't join this game. Check the invitation link or try again.
        <v-btn text small color="error" @click="retry">Try again</v-btn>
        <v-btn text small to="/">Go home</v-btn>
      </v-alert>
      <v-alert v-if="error" type="error" text dismissible role="alert" @input="error = ''">{{ error }}</v-alert>

      <section v-if="!ready && !joinFailed && !error" class="joining-state" aria-live="polite" aria-busy="true">
        <v-progress-circular indeterminate color="primary" size="28"></v-progress-circular>
        <p class="muted">Joining your game…</p>
      </section>
      <v-btn v-if="!ready && error" color="primary" depressed @click="retry">Try again</v-btn>

      <div v-if="ready" class="game-workspace" :class="{ 'game-workspace-finished': state.isFinished() }">
        <section class="table-panel" aria-labelledby="round-title">
          <div class="round-heading">
            <div>
              <h2 id="round-title">{{ roundTitle }}</h2>
              <p class="round-description muted">{{ roundDescription }}</p>
            </div>
            <span class="round-badge" :class="{ 'round-badge-finished': state.isFinished() }">
              <span class="status-dot" aria-hidden="true"></span>
              {{ state.isFinished() ? 'Revealed' : 'Voting open' }}
            </span>
          </div>

          <div class="players-scroll" role="region" aria-label="Players and votes" tabindex="0">
            <div class="players-table">
              <CardOnTable
                v-for="(player, index) in players"
                :key="index"
                :name="player.name"
                :card="player.voted_card"
                :confidence="player.confidence"
                :finished="state.isFinished()"
              ></CardOnTable>
            </div>
          </div>

          <div class="round-actions">
            <p class="vote-progress" role="status" aria-live="polite">{{ votedCount }} of {{ players.length }} players voted</p>
            <v-btn
              v-if="state.canReveal()"
              color="primary"
              large
              depressed
              :loading="pendingAction === 'reveal'"
              :disabled="!connected"
              :aria-disabled="busy || !connected"
              @click="runAction('reveal')"
            >Reveal cards</v-btn>
            <v-btn
              v-else-if="state.canRestart()"
              color="primary"
              large
              depressed
              :loading="pendingAction === 'restart'"
              :disabled="!connected"
              :aria-disabled="busy || !connected"
              @click="runAction('restart')"
            >Start next round</v-btn>
            <p v-else class="muted round-permission">{{ state.isFinished() ? 'Waiting for the next round.' : 'Pick a card while the team votes.' }}</p>
          </div>
        </section>

        <section v-if="state.isRunning()" class="deck-panel" aria-labelledby="deck-title">
          <div class="deck-heading">
            <h2 id="deck-title">Your estimate</h2>
            <span class="vote-feedback" role="status" aria-live="polite">
              <template v-if="pendingAction === 'vote'">Saving your vote…</template>
              <template v-else-if="!connected">Waiting for connection</template>
              <template v-else-if="state.voted()">✓ Vote saved</template>
              <template v-else>Pick a card</template>
            </span>
          </div>

          <div class="voting-deck" role="group" aria-label="Choose your estimate">
            <Card
              v-for="card in state.getCards()"
              :key="card"
              :value="card"
              :confidence="state.confidence"
              :active="state.isActive(card)"
              :disabled="!connected"
              :busy="busy"
              @click="runAction('vote', card)"
            ></Card>
          </div>

          <div class="confidence-control">
            <span id="confidence-label" class="muted">How confident are you?</span>
            <div class="confidence-options" role="group" aria-labelledby="confidence-label">
              <button
                v-for="option in confidenceOptions"
                :key="option.value"
                type="button"
                :aria-pressed="state.voted() && state.confidence === option.value"
                :class="{ 'confidence-selected': state.voted() && state.confidence === option.value }"
                :disabled="!connected || !state.voted()"
                :aria-disabled="busy || !connected || !state.voted()"
                @click="runAction('changeConfidence', option.value)"
              >{{ option.label }}</button>
            </div>
          </div>
        </section>

        <section v-else class="finished-panel" aria-labelledby="finished-title">
          <h2 id="finished-title">Round complete</h2>
          <p class="muted">Start the next round when you're ready to vote again.</p>
        </section>
      </div>
    </main>

    <v-snackbar v-model="showCopyNotice" color="#242424" :timeout="3000" role="status">
      Invitation link copied.
      <template v-slot:action="{ attrs }">
        <v-btn text v-bind="attrs" aria-label="Dismiss notification" @click="showCopyNotice = false">Close</v-btn>
      </template>
    </v-snackbar>

    <UserNameDialog ref="userNameDialog"></UserNameDialog>
    <GameSettingsDialog ref="gameSettingsDialog"></GameSettingsDialog>
    <InviteDialog ref="inviteDialog"></InviteDialog>
  </div>
</template>

<script>
import State from "@/models/state"
import game from "@/models/game"
import user from "@/models/user"
import notifier from "@/notifier/notifier"
import Card from "@/components/Card"
import CardOnTable from "@/components/CardOnTable"
import UserNameDialog from "@/components/UserNameDialog"
import InviteDialog from "@/components/InviteDialog"
import GameSettingsDialog from "@/components/GameSettingsDialog"

export default {
  name: 'Game',

  components: {
    UserNameDialog,
    CardOnTable,
    Card,
    InviteDialog,
    GameSettingsDialog,
  },

  data: () => ({
    user,
    state: null,
    notifier,
    error: '',
    pendingAction: '',
    showCopyNotice: false,
    confidenceOptions: [
      {value: 'low', label: 'Not sure'},
      {value: 'normal', label: 'Default'},
      {value: 'high', label: 'Very confident'},
    ],
  }),

  computed: {
    ready() {
      return this.state && (this.state.isRunning() || this.state.isFinished())
    },

    connected() {
      return this.notifier.status === this.notifier.STATUS_CONNECTED &&
        this.notifier.listenStatus === this.notifier.STATUS_JOINED
    },

    joinFailed() {
      return this.notifier.listenStatus === this.notifier.STATUS_JOIN_FAILED
    },

    busy() {
      return this.pendingAction !== ''
    },

    players() {
      return this.state.getPlayers()
    },

    votedCount() {
      return this.players.filter(player => !!player.voted_card).length
    },

    roundTitle() {
      if (this.state.isFinished()) {
        return 'Cards revealed'
      }

      return this.players.length > 0 && this.votedCount === this.players.length
        ? 'Everyone has voted'
        : 'What’s your estimate?'
    },

    roundDescription() {
      if (this.state.isFinished()) {
        return 'Compare the estimates and discuss any differences.'
      }

      return this.votedCount === this.players.length
        ? 'Ready to reveal the cards?'
        : 'Votes stay hidden until you reveal the cards.'
    },

    ticketLink() {
      try {
        const url = new URL(this.state.ticket_url)
        return ['http:', 'https:'].includes(url.protocol) ? url.href : ''
      } catch (error) {
        return ''
      }
    },
  },

  async mounted() {
    try {
      if (!user.identified()) {
        user.name = await this.$refs.userNameDialog.open(user.name)
      }

      await user.authenticate()
      this.state = new State(this.$route.params.id)
    } catch (error) {
      this.error = "We couldn't connect to your game. Check your connection and try again."
    }
  },

  beforeDestroy() {
    if (this.state) {
      this.state.stopUpdates()
    }
  },

  methods: {
    retry() {
      window.location.reload()
    },

    async runAction(action, value) {
      if (this.busy || !this.connected) {
        return
      }

      this.pendingAction = action
      this.error = ''

      try {
        await this.state[action](value)
      } catch (error) {
        this.error = "That change couldn't be saved. Please try again."
      } finally {
        this.pendingAction = ''
      }
    },

    async changeName() {
      const name = await this.$refs.userNameDialog.openModify(user.name)
      if (!name) {
        return
      }

      try {
        await user.update(name)
      } catch (error) {
        this.error = "Your name couldn't be saved. Please try again."
      }
    },

    async changeGameParams() {
      const [ok, name, url] = await this.$refs.gameSettingsDialog.openModify(this.state.name, this.state.ticket_url)
      if (!ok) {
        return
      }

      try {
        await game.update(this.state.id, name, url)
      } catch (error) {
        this.error = "Your game settings couldn't be saved. Please try again."
      }
    },

    async copyLink() {
      try {
        await navigator.clipboard.writeText(location.href)
        this.showCopyNotice = true
      } catch (error) {
        this.showCopyNotice = await this.$refs.inviteDialog.open(location.href)
      }
    },
  },
}
</script>

<style scoped>
.game-page {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.game-header {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  gap: 24px;
  min-height: 72px;
  padding: 16px 32px;
  border-bottom: 1px solid var(--poker-border);
  background: var(--poker-surface);
}

.session-title {
  min-width: 0;
  padding-left: 24px;
  border-left: 1px solid var(--poker-border);
}

.session-title h1 {
  display: -webkit-box;
  overflow: hidden;
  font-size: 1rem;
  line-height: 1.5;
  overflow-wrap: anywhere;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.session-title a {
  color: var(--poker-text);
  text-decoration: none;
}

.session-title a:hover {
  color: var(--poker-accent-ink);
  text-decoration: underline;
}

.session-deck {
  font-size: 0.75rem;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
  flex-shrink: 0;
}

.game-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
  width: 100%;
  max-width: 1280px;
  padding: 24px 32px;
  margin: 0 auto;
}

.game-content > .v-alert {
  flex-shrink: 0;
  margin-bottom: 0;
}

.game-workspace {
  display: grid;
  gap: 20px;
  min-width: 0;
}

.joining-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  padding: 80px 0;
}

.table-panel {
  display: grid;
  grid-template-rows: auto minmax(0, 1fr) auto;
  gap: 16px;
  min-width: 0;
  min-height: 0;
  padding: 24px;
  background: var(--poker-surface);
  border: 1px solid var(--poker-border);
  border-radius: var(--poker-radius);
}

.round-heading {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.round-heading h2 {
  font-size: 1.5rem;
  letter-spacing: -0.035em;
  line-height: 1.3;
}

.round-description {
  margin: 8px 0 0;
  font-size: 0.9rem;
  line-height: 1.6;
}

.round-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
  padding: 7px 12px;
  border-radius: 20px;
  background: #f0f0f0;
  color: var(--poker-muted);
  font-size: 0.75rem;
  font-weight: 600;
}

.round-badge-finished {
  background: var(--poker-accent-soft);
  color: var(--poker-accent-ink);
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
}

.players-scroll {
  min-height: 0;
  max-height: 320px;
  overflow: auto;
  overscroll-behavior: contain;
}

.players-scroll:focus-visible {
  outline: 3px solid var(--poker-accent-ink);
  outline-offset: 2px;
}

.players-table {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  align-content: center;
  justify-content: center;
  gap: 24px 12px;
  min-height: 100%;
  padding: 8px 4px;
}

.round-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: center;
  gap: 12px 20px;
  min-height: 44px;
}

.vote-progress {
  margin: 0;
  color: var(--poker-muted);
  font-size: 0.8rem;
}

.round-permission {
  margin: 8px 0 0;
  text-align: center;
  font-size: 0.9rem;
}

.deck-panel {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 16px;
  min-width: 0;
  padding: 24px;
  background: var(--poker-surface);
  border: 1px solid var(--poker-border);
  border-radius: var(--poker-radius);
}

.deck-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.deck-heading h2 {
  font-size: 1rem;
  font-weight: 600;
}

.vote-feedback {
  color: var(--poker-accent-ink);
  font-size: 0.8rem;
}

.voting-deck {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: center;
  gap: 12px 10px;
  padding: 10px 0 6px;
}

.confidence-control {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 8px;
  font-size: 0.8rem;
}

.confidence-options {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 4px;
  padding: 4px;
  border: 1px solid var(--poker-border);
  border-radius: 10px;
}

.confidence-options button {
  min-height: 44px;
  padding: 6px 12px;
  border-radius: 7px;
  color: var(--poker-muted);
  font-size: 0.8rem;
  font-weight: 500;
}

.confidence-options button:hover:not([aria-disabled="true"]),
.confidence-options .confidence-selected {
  background: var(--poker-accent-soft);
  color: var(--poker-accent-ink);
}

.confidence-options button[aria-disabled="true"] {
  opacity: 0.6;
  cursor: wait;
}

.confidence-options button:disabled {
  cursor: default;
  opacity: 0.5;
}

.finished-panel {
  padding: 12px;
  text-align: center;
}

.finished-panel h2 {
  font-size: 1.1rem;
  font-weight: 600;
}

.finished-panel p {
  max-width: 540px;
  margin: 12px auto 0;
  font-size: 0.9rem;
  line-height: 1.7;
}

/* Keep voting visible, with only the player area scrolling. */
@media (min-width: 700px) and (min-height: 600px), (min-width: 900px) and (min-height: 480px) and (max-height: 649px) {
  .game-page {
    height: 100vh;
    height: 100dvh;
  }

  .game-content {
    flex: 1;
    min-height: 0;
    padding: 16px 24px;
  }

  .game-workspace {
    flex: 1;
    min-height: 0;
    grid-template-rows: minmax(0, 1fr) auto;
    gap: 16px;
  }

  .table-panel {
    gap: 12px;
    padding: 20px;
  }

  .deck-panel {
    gap: 12px;
    padding: 16px 20px;
  }

  .players-scroll {
    max-height: none;
  }

  .finished-panel p {
    margin-top: 8px;
  }
}

@media (min-width: 700px) and (max-height: 740px) {
  .game-workspace {
    --vote-card-width: 48px;
    --vote-card-height: 80px;
  }

  .table-panel, .deck-panel {
    padding: 16px;
  }

  .round-description {
    display: none;
  }

  .players-table {
    padding-top: 0;
    padding-bottom: 0;
  }

  .voting-deck {
    column-gap: 8px;
  }
}

@media (min-width: 700px) and (max-width: 959px) {
  .game-workspace {
    --vote-card-width: 44px;
    --vote-card-height: 74px;
  }

  .voting-deck {
    column-gap: 6px;
  }
}

/* A short, wide window is the exception to the board-above-deck layout. */
@media (min-width: 900px) and (min-height: 480px) and (max-height: 649px) {
  .game-workspace {
    grid-template-columns: minmax(0, 1fr) minmax(360px, 0.85fr);
    grid-template-rows: minmax(0, 1fr);
    gap: 20px;
  }

  .game-workspace-finished {
    grid-template-columns: minmax(0, 1fr);
    grid-template-rows: minmax(0, 1fr) auto;
  }
}

@media (max-width: 899px) {
  .game-header {
    gap: 20px;
    padding: 16px 24px;
  }

  .game-brand-name {
    display: none;
  }

  .session-title {
    padding-left: 20px;
  }
}

@media (max-width: 600px) {
  .game-header {
    flex-wrap: wrap;
    gap: 12px;
    padding: 16px;
  }

  .session-title {
    order: 3;
    flex-basis: 100%;
    padding: 12px 0 0;
    border-left: 0;
    border-top: 1px solid var(--poker-border);
  }

  .game-content {
    padding: 16px;
  }

  .table-panel, .deck-panel {
    padding: 20px 16px;
  }

  .round-heading {
    flex-direction: column-reverse;
    gap: 12px;
  }

  .round-heading h2 {
    font-size: 1.5rem;
  }

  .players-table {
    gap: 20px 8px;
  }

  .deck-heading {
    align-items: flex-start;
    flex-direction: column;
    gap: 8px;
  }

  .voting-deck {
    gap: 10px;
  }

  .confidence-control {
    flex-direction: column;
  }

  .confidence-options button {
    min-height: 44px;
    padding: 6px 8px;
    font-size: 0.75rem;
  }

  .confidence-options {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    width: 100%;
    max-width: 380px;
  }

  .invite-label {
    display: none;
  }

  .finished-panel {
    padding: 12px 8px;
  }
}
</style>
