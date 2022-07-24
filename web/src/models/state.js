import game from "@/models/game";
import notifier from "@/notifier/notifier";

const stateRunning = 'started'
const stateFinished = 'finished'

export default class {
    id;
    name;
    ticket_url;
    state;
    cards_deck;
    players;
    voted_card;
    can_reveal;
    confidence;

    constructor(id) {
        this.id = id
        notifier.listenGame(id, state => this.updateState(state))
    }

    getCards() {
        if (this.cards_deck) {
            return this.cards_deck.cards
        }
        return []
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
        return this.can_reveal && this.state === stateRunning
    }

    canRestart() {
        return this.can_reveal && this.state === stateFinished
    }

    isActive(card) {
        return this.voted_card === card
    }

    voted() {
        return this.voted_card !== ""
    }

    async reveal() {
        await game.reveal(this.id)
    }

    async restart() {
        await game.restart(this.id)
    }

    async vote(card) {
        if (this.voted_card === card) {
            this.voted_card = ""
            await game.unVote(this.id)
        } else {
            card = encodeURIComponent(card)
            this.voted_card = card
            this.confidence = "normal"
            await game.vote(this.id, card, this.confidence)
        }
    }

    async changeConfidence(confidence) {
        if (!this.voted_card || this.voted_card === "") {
            return
        }
        await game.vote(this.id, this.voted_card, confidence)
    }

    stopUpdates() {
        notifier.leaveGame()
    }

    updateState(state) {
        for (const attribute in state) {
            this[attribute] = state[attribute];
        }
        this.players && this.players.sort(comparePlayers)
    }
}

function comparePlayers(a, b) {
    if (a.name < b.name) {
        return -1;
    }
    if (a.name > b.name) {
        return 1;
    }
    return 0;
}
