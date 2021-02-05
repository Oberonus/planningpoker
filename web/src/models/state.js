import game from "@/models/game";

const stateRunning = 'started'
const stateFinished = 'finished'

let instance = null, callback;

export default class State{
    workerID;
    id;
    State;
    Cards;
    Players;
    VotedCard;
    CanReveal;

    static getInstance = (func) => {
        if (instance) {
            func(instance);
        } else {
            callback = func;
        }
    }
    static setInstance = (value) => {
        instance = value;
        callback && callback(instance);
    }

    constructor(id) {
        this.id = id;
        State.setInstance(this);
    }

    cards() {
        return this.Cards
    }

    isRunning() {
        return this.State === stateRunning
    }

    isFinished() {
        return this.State === stateFinished
    }

    players() {
        return this.Players
    }

    canReveal() {
        return this.CanReveal && this.State === stateRunning
    }

    canRestart() {
        return this.CanReveal && this.State === stateFinished
    }

    isActive(card) {
        return this.VotedCard === card
    }

    async reveal() {
        await game.reveal(this.id)
        await this.update()
    }

    async restart() {
        await game.restart(this.id)
        await this.update()
    }

    async vote(card) {
        if (this.VotedCard === card) {
            this.VotedCard = ""
            await game.unVote(this.id)
            await this.update()
        } else {
            card = encodeURIComponent(card)
            this.VotedCard = card
            await game.vote(this.id, card)
            await this.update()
        }
    }

    updatePeriodically() {
        this.workerID = setInterval(() => {
            this.updatePlayers()
                .catch(e => {
                    console.error(e);
                    this.stopUpdates();
                    throw new Error('Error');
                })
        }, 500)
    }

    stopUpdates() {
        clearInterval(this.workerID)
    }

    update() {
        return game
            .state(this.id)
            .then((state) => {
                for (const key in state) {
                    this[key] = state[key];
                }
                this.Players && this.Players.sort(comparePlayers)
            })
    }

    updatePlayers() {
        return game.state(this.id)
            .then(state => {
                this.State = state.State
                this.CanReveal = state.CanReveal
                this.Players = state.Players
                this.Players && this.Players.sort(comparePlayers)
                if (this.State === stateFinished) {
                    this.VotedCard = ""
                }
            })
    }
}

function comparePlayers(a, b) {
    if (a.Name < b.Name) {
        return -1;
    }
    if (a.Name > b.Name) {
        return 1;
    }
    return 0;
}
