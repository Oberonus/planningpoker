<template>
  <v-container fill-height fluid v-if="state">
      <v-row align="center" justify="center" style="min-height: 100px;"/>

      <v-row align="center" justify="center">
        <card-on-table v-for="player in state.players()" v-bind:key="player.Name" :name="player.Name"
                       :card="player.VotedCard"></card-on-table>
      </v-row>

      <v-row align="center" justify="center" style="min-height: 100px;">
        <v-btn v-if="state.canReveal()" @click="state.reveal()">Show cards</v-btn>
        <v-btn v-else-if="state.canRestart()" @click="state.restart()">New voting</v-btn>
        <span v-if="state.isRunning() && !state.canReveal()">Please pick your cards</span>
      </v-row>

      <v-row v-if="isRunning"
          align="center"
          justify="center"
          style="min-height: 8rem;"
      >
          <v-col cols="1"
                 class="pa-2"
                 v-for="card in cards"
                 :key="card"
                 style="min-width: 4rem; max-width: 4.5rem"
          >
            <Card :value="card"
                  :active="state.isActive(card)"
                  @select="onSelectCard(card)"
            />
          </v-col>
      </v-row>

      <v-row v-else-if="state.isFinished()"
           align="center"
           justify="center"
           style="min-height: 8rem;"
      >
          <h1>Voting finished</h1>
      </v-row>

    </v-container>
</template>

<script>

import Card from "@/components/Card";
import CardOnTable from "@/components/CardOnTable";
import State from "@/models/state";
import game from "@/models/game";

export default {
  name: 'Game',

  components: {
    CardOnTable,
    Card,
  },

  props: {
    user: {
      required: true,
    },
  },

  data: () => ({
      state: null,
      invitationOpacity: 0,
      gameID: null,
  }),

  watch:{
    "user.identified": function(value) {
      value && this.startGame();
    },
  },

  mounted() {
      this.gameID = this.$route.params.id;
      this.startGame();
  },

  async beforeDestroy() {
    this.state && this.state.stopUpdates()
  },

  computed: {
    cards() {
      return this.state.cards() || [];
    },

    isRunning() {
      return this.state && this.state.isRunning();
    }
  },

  methods: {
    startGame(){
      this.user.authenticate().then(() => {
        game.join(this.gameID).then(() => {
          this.state = new State(this.gameID);
          this.state.updatePeriodically();
          this.state.update();
        })
      }).catch(e => {
        console.warn(e);
      })
    },

    onSelectCard(card){
      //TODO: animate playground card
      this.state.vote(card);
    },

  }
}
</script>

<style scoped>
</style>
