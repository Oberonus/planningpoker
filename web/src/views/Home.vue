<template>
  <v-container fill-height fluid>
    <v-row align="center" justify="center">
      <v-container>
        <v-row align="center" justify="center" style="margin-bottom: 10pt;">
          <div style="display: block" v-if="user.identified"><h1>Hi {{ user.name }}!</h1></div>
        </v-row>
        <v-row align="center" justify="center">
          <v-btn @click="startGame" x-large color="#23D2AA" dark> Start new game</v-btn>
        </v-row>
      </v-container>
    </v-row>
  </v-container>
</template>

<script>

import game from "@/models/game"

export default {
  props:{
    user : {
      required: true,
    },
  },

  methods: {
    async startGame() {
      const id = await game.create();
      return this.$router.push({
        name: 'Games',
        params: { id },
      })
    }
  }
}
</script>
