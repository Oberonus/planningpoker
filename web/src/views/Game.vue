<template>
  <v-scroll-y-transition>
    <v-container fill-height fluid v-if="state">
      <v-row no-gutters align="stretch" class="fill-height">

        <v-col cols="12" class="row align-end justify-center no-gutters">
          <card-on-table v-for="player in state.players()"
                         :key="player.Name" :name="player.Name"
                         :card="player.VotedCard"/>
        </v-col>

        <v-col cols="12" class="row align-center justify-center mt-0 no-gutters">
            <v-btn v-if="actionsPossible"
                   width="120"
                   @click="onClick"
            >
              <span v-text="actionText" />
            </v-btn>

            <span v-if="isRunning && !state.canReveal()"
                  v-text="'Please pick your cards'"
            />
        </v-col>

        <v-col cols="12"
               class="row justify-center align-center mt-0 no-gutters"
               style="min-height: 10rem"
        >
          <v-row class="fill-height" align="center" justify="center">
            <v-scale-transition hide-on-leave leave-absolute>
              <v-row v-if="isRunning && !isFinished" no-gutters class="fill-height" align="center" justify="center">
                <v-col class="px-2 py-0"
                       v-for="card in cards"
                       :key="card"
                       style="min-width: 4rem; max-width: 4.5rem"
                >
                  <card :value="card"
                        :active="state.isActive(card)"
                        @select="onSelectCard(card)"
                  />
                </v-col>
              </v-row>
            </v-scale-transition>

            <v-scale-transition hide-on-leave leave-absolute>
              <h1 v-show="isFinished" v-text="'Voting finished'"/>
            </v-scale-transition>
          </v-row>

        </v-col>

      </v-row>
    </v-container>
  </v-scroll-y-transition>
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

  beforeDestroy() {
    this.state && this.state.stopUpdates();
  },

  computed: {
    cards() {
      return this.state.cards() || [];
    },

    isRunning() {
      return this.state && this.state.isRunning();
    },

    isFinished(){
      return this.state && this.state.isFinished();
    },

    actionsPossible(){
      return this.state.canReveal() || this.state.canRestart();
    },

    actionText(){
      return this.state.canRestart() ? 'New voting' : 'Show cards';
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

    onClick(){
      this.state.canReveal() && this.state.reveal();
      this.state.canRestart() && this.state.restart();
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
