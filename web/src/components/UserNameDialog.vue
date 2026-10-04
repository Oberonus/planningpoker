<template>
  <v-dialog v-model="show" max-width="440px" :persistent="!modifying">
    <v-card class="dialog-card" role="dialog" aria-modal="true" aria-labelledby="name-dialog-title">
      <v-form ref="nameForm" @submit.prevent="save" @keydown.esc.native.stop.prevent="modifying && cancel()">
        <v-card-title id="name-dialog-title" role="heading" aria-level="2">{{ modifying ? 'Your name' : 'Join game' }}</v-card-title>
        <v-card-text>
          <p class="dialog-description">{{ modifying ? 'Choose the name your teammates see.' : 'What should your teammates call you?' }}</p>
          <v-text-field
            ref="nameInput"
            v-model="name"
            label="Your name"
            autocomplete="nickname"
            :rules="nameRules"
            outlined
            dense
          ></v-text-field>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn v-if="modifying" text @click="cancel">Cancel</v-btn>
          <v-btn type="submit" color="primary" depressed>{{ modifying ? 'Save name' : 'Join game' }}</v-btn>
        </v-card-actions>
      </v-form>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  data: () => ({
    modifying: false,
    show: false,
    name: '',
    resolve: null,
    nameRules: [value => !!value.trim() || 'Enter your name to join the game.'],
  }),

  watch: {
    show(visible) {
      if (!visible && this.resolve) {
        this.resolve('')
        this.resolve = null
      }
    },
  },

  methods: {
    open(name) {
      return this.showDialog(name, false)
    },

    openModify(name) {
      return this.showDialog(name, true)
    },

    showDialog(name, modifying) {
      this.modifying = modifying
      this.name = name || ''
      this.show = true

      this.$nextTick(() => {
        this.$refs.nameForm.resetValidation()
        // Vuetify focuses its lazy dialog content after two render ticks.
        requestAnimationFrame(() => {
          if (this.show) {
            this.$refs.nameInput.focus()
          }
        })
      })

      return new Promise(resolve => {
        this.resolve = resolve
      })
    },

    save() {
      if (!this.$refs.nameForm.validate()) {
        return
      }

      this.resolve(this.name.trim())
      this.resolve = null
      this.show = false
    },

    cancel() {
      this.show = false
    },
  },
}
</script>
