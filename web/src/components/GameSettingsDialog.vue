<template>
  <v-dialog v-model="show" max-width="500px" persistent>
    <v-card>
      <v-card-title>
        <span class="headline">Setup your game</span>
      </v-card-title>
      <v-divider></v-divider>
      <v-card-text style="margin-top: 20px;">
        <v-row>
          <v-col cols="12">
            <v-text-field v-model="name" label="Name of the game (optional)"></v-text-field>
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12">
            <v-text-field v-model="url" label="Ticket URL (optional)"></v-text-field>
          </v-col>
        </v-row>
        <v-row v-if="deck!==null" align="center">
          <v-col cols="12">
            <v-select
                v-model="deck"
                :items="decks"
                :item-text="item => formatDeck(item)"
                label="Card deck"
                return-object
            ></v-select>
          </v-col>
        </v-row>
      </v-card-text>
      <v-divider></v-divider>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary" text @click="save">OK</v-btn>
        <v-btn color="primary" text @click="cancel">Cancel</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import game from "@/models/game"

export default {
  data: () => {
    return {
      modeCreate: 'create',
      modeModify: 'modify',
      mode: null,
      show: false,
      name: '',
      url: '',
      resolve: null,
      deck: null,
      decks: game.decks,
    }
  },

  methods: {
    async open() {
      this.mode = this.modeCreate
      this.name = ""
      this.url = ""
      this.show = true
      this.deck = this.decks[0]

      return new Promise((resolve) => {
        this.resolve = resolve
      })
    },

    async openModify(name, url) {
      this.mode = this.modeModify
      this.name = name
      this.url = url
      this.show = true
      this.deck = null

      return new Promise((resolve) => {
        this.resolve = resolve
      })
    },

    save() {
      this.show = false
      this.resolve([true, this.name, this.url, this.deck])
    },

    cancel() {
      this.show = false
      this.resolve([false, "", "", null])
    },

    formatDeck(deck) {
      return deck.name + " (" + deck.types.join(", ") + ")"
    }
  },
}
</script>