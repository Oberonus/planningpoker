import game from "@/models/game"
import notifier from "@/notifier/notifier"

const stateRunning = 'started'
const stateFinished = 'finished'

export default class {
    id = '';
    name = '';
    ticket_url = '';
    state = '';
    cards_deck = {name: '', cards: []};
    players = [];
    voted_card = '';
    can_reveal = false;
    confidence = 'normal';

    constructor(id) {
        this.id = id
        notifier.listenGame(id, state => this.updateState(state))
    }

    getCards() {
        return this.cards_deck.cards
    }

    isRunning() {
        return this.state === stateRunning
    }

    isFinished() {
        return this.state === stateFinished
    }

    getPlayers() {
        return this.players
    }

    canReveal() {
        return this.can_reveal && this.isRunning()
    }

    canRestart() {
        return this.can_reveal && this.isFinished()
    }

    isActive(card) {
        return this.voted_card === card
    }

    voted() {
        return this.voted_card !== ''
    }

    reveal() {
        return game.reveal(this.id)
    }

    restart() {
        return game.restart(this.id)
    }

    vote(card) {
        if (this.voted_card === card) {
            return game.unVote(this.id)
        }

        // Selection and saved feedback follow the server's personalized state update.
        return game.vote(this.id, card, 'normal')
    }

    changeConfidence(confidence) {
        if (!this.voted()) {
            return
        }

        return game.vote(this.id, this.voted_card, confidence)
    }

    stopUpdates() {
        notifier.leaveGame()
    }

    updateState(state) {
        for (const attribute in state) {
            this[attribute] = state[attribute]
        }

        this.players.sort(comparePlayers)
    }
}

function comparePlayers(first, second) {
    return first.name.localeCompare(second.name)
}
