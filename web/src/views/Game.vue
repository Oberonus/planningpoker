<template>
  <v-container fill-height fluid>
    <v-app-bar absolute>
      <v-toolbar-title v-if="state" style="font-weight: 700">
        <v-btn icon href="/">
          <v-icon>mdi-home</v-icon>
        </v-btn>
        <v-btn v-if="state.TicketURL" :depressed="true" :ripple="false"
               :href="state.TicketURL" target="_blank">
          <span>{{ state.Name ? state.Name : state.TicketURL }}</span>
        </v-btn>
        <v-btn v-if="!state.TicketURL && state.Name" :depressed="true" :ripple="false" target="_blank">
          <span>{{ state.Name }}</span>
        </v-btn>
        <v-btn icon @click="changeGameParams">
          <v-icon>mdi-cog</v-icon>
        </v-btn>
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <v-btn @click="copyLink">
        <v-icon>mdi-content-copy</v-icon>
        <span style="margin-left: 10px;">Invite players</span>
      </v-btn>
      <v-btn style="margin-left: 10px;" icon @click="changeName">
        <v-icon>mdi-account-edit</v-icon>
      </v-btn>
    </v-app-bar>

    <v-banner dark color="rgb(63,29,203)"
              :style="'position:absolute; top:80px; left:10%; right: 10%; text-align: center; transition: opacity 0.3s ease-in-out; opacity: '+invitationOpacity+';'"
              elevation="3">Invitation link copied to clipboard!
    </v-banner>

    <v-row align="center" justify="center" v-if="state" style="min-height: 100px;">
    </v-row>

    <v-row align="center" justify="center" v-if="state">
      <card-on-table v-for="player in state.players()" v-bind:key="player.Name" :name="player.Name"
                     :card="player.VotedCard"></card-on-table>
    </v-row>

    <v-row align="center" justify="center" v-if="state" style="min-height: 100px;">
      <v-btn v-if="state.canReveal()" @click="state.reveal()">Show cards</v-btn>
      <v-btn v-else-if="state.canRestart()" @click="state.restart()">New voting</v-btn>
      <div v-if="state.isRunning() && !state.canReveal()">Please pick your cards</div>
    </v-row>

    <v-row align="center" justify="center" v-if="state" v-show="state.isRunning()" style="min-height: 100px;">
      <Card v-for="card in state.cards()"
            v-bind:key="card"
            @click="state.vote(card)"
            :active="state.isActive(card)"
            :value="card">
      </Card>
    </v-row>

    <v-row align="center" justify="center" v-if="state" v-show="state.isFinished()" style="min-height: 100px;">
      <h1>Voting finished</h1>
    </v-row>
    <UserNameDialog ref="userNameDialog"></UserNameDialog>
    <GameSettingsDialog ref="gameSettingsDialog"></GameSettingsDialog>
    <InviteDialog ref="inviteDialog"></InviteDialog>
  </v-container>
</template>

<script>

import State from "@/models/state";
import game from "@/models/game";
import user from "@/models/user";
import Card from "@/components/Card";
import UserNameDialog from "@/components/UserNameDialog";
import CardOnTable from "@/components/CardOnTable";
import InviteDialog from "@/components/InviteDialog";
import GameSettingsDialog from "@/components/GameSettingsDialog"

export default {
  name: 'Game',

  components: {
    UserNameDialog,
    CardOnTable,
    Card,
    InviteDialog,
    GameSettingsDialog
  },

  data: () => {
    return {
      user: user,
      state: null,
      invitationOpacity: 0,
    }
  },

  async mounted() {
    try {
      if (!user.identified()) {
        user.name = await this.$refs.userNameDialog.open(user.name)
      }
      await user.authenticate()

      const gameID = this.$route.params.id
      await game.join(gameID)
      this.state = new State(gameID)
    } catch (e) {
      console.log(e)
      await this.$router.push({name: 'Home'})
    }
  },

  async beforeDestroy() {
    this.state && this.state.stopUpdates()
  },

  methods: {
    async changeName() {
      const name = await this.$refs.userNameDialog.openModify(user.name)
      if (name !== "") {
        await user.update(name)
      }
    },

    async changeGameParams() {
      const [ok, name, url] = await this.$refs.gameSettingsDialog.openModify(this.state.Name, this.state.TicketURL)
      if (!ok) {
        return
      }
      await game.update(this.state.id, name, url)
    },

    async copyLink() {
      await this.$refs.inviteDialog.open(location.href)
      this.invitationOpacity = 1
      setTimeout(() => {
        this.invitationOpacity = 0
      }, 3000)
    }
  }
}
</script>
