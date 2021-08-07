import game from "@/models/game";

const stateRunning = 'started'
const stateFinished = 'finished'

export default class {
    workerID;

    id;
    State;
    Cards;
    Players;
    VotedCard;
    CanReveal;
    ChangeID;

    constructor(id) {
        this.id = id
        this.updatePeriodically()
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
    }

    async restart() {
        await game.restart(this.id)
    }

    async vote(card) {
        if (this.VotedCard === card) {
            this.VotedCard = ""
            await game.unVote(this.id)
        } else {
            card = encodeURIComponent(card)
            this.VotedCard = card
            await game.vote(this.id, card)
        }
    }

    // long polling operation, will throttle in case of any exception
    async updatePeriodically() {
        try {
            await this.update()
            setTimeout(async () => {
                await this.updatePeriodically()
            }, 0)
        } catch (e) {
            setTimeout(async () => {
                await this.updatePeriodically()
            }, 1000)
            throw e
        }
    }

    stopUpdates() {
        clearInterval(this.workerID)
    }

    async update() {
        const state = await game.state(this.id, this.ChangeID)
        for (const attribute in state) {
            this[attribute] = state[attribute];
        }
        this.Players && this.Players.sort(comparePlayers)
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
