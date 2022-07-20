<template>
  <v-container fill-height fluid>
    <v-row align="center" justify="center">
      <v-container>
        <v-row align="center" justify="center" style="margin-bottom: 10pt;">
          <div style="display: block" v-if="user.identified()"><h1>Hi {{ user.name }}!</h1></div>
        </v-row>
        <v-row align="center" justify="center">
          <v-btn @click="startGame" x-large color="#23D2AA" dark> Start new game</v-btn>
        </v-row>
      </v-container>
    </v-row>
    <UserNameDialog ref="userNameDialog"></UserNameDialog>
    <NewGameDialog ref="newGameDialog"></NewGameDialog>
  </v-container>
</template>

<script>

import game from "@/models/game"
import user from "@/models/user"
import UserNameDialog from "@/components/UserNameDialog";
import NewGameDialog from "@/components/GameSettingsDialog"

export default {
  data: () => {
    return {
      user,
    }
  },

  components: {
    UserNameDialog,
    NewGameDialog
  },

  methods: {
    async startGame() {
      if (!user.identified()) {
        user.name = await this.$refs.userNameDialog.open(user.name)
      }
      await user.authenticate()

      const [result, gameName, gameURL, deck] = await this.$refs.newGameDialog.open()
      if (!result) {
        return
      }

      this.$router.push({
        name: 'Games',
        params: {id: await game.create(gameName, gameURL, deck)},
      })
    }
  }
}
</script>
