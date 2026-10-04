<template>
  <v-dialog v-model="show" max-width="480px">
    <v-card class="dialog-card" role="dialog" aria-modal="true" aria-labelledby="settings-dialog-title">
      <v-form ref="settingsForm" @submit.prevent="save" @keydown.esc.native.stop.prevent="cancel">
        <v-card-title id="settings-dialog-title" role="heading" aria-level="2">Game settings</v-card-title>
        <v-card-text>
          <p class="dialog-description">Update the game name or ticket link.</p>
          <v-text-field
            ref="gameNameInput"
            v-model="name"
            label="Game name (optional)"
            placeholder="e.g. Sprint planning"
            outlined
            dense
          ></v-text-field>
          <v-text-field
            v-model="url"
            label="Ticket URL (optional)"
            placeholder="https://"
            :rules="urlRules"
            outlined
            dense
          ></v-text-field>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="cancel">Cancel</v-btn>
          <v-btn type="submit" color="primary" depressed>Save settings</v-btn>
        </v-card-actions>
      </v-form>
    </v-card>
  </v-dialog>
</template>

<script>
import { ticketURLRules } from "@/models/validation"

export default {
  data: () => ({
    show: false,
    name: '',
    url: '',
    resolve: null,
    urlRules: ticketURLRules,
  }),

  watch: {
    show(visible) {
      if (!visible && this.resolve) {
        this.resolve([false, '', ''])
        this.resolve = null
      }
    },
  },

  methods: {
    openModify(name, url) {
      this.name = name || ''
      this.url = url || ''
      this.show = true

      this.$nextTick(() => {
        this.$refs.settingsForm.resetValidation()
        // Wait for Vuetify's lazy dialog to finish moving focus.
        requestAnimationFrame(() => {
          if (this.show) {
            this.$refs.gameNameInput.focus()
          }
        })
      })

      return new Promise(resolve => {
        this.resolve = resolve
      })
    },

    save() {
      if (!this.$refs.settingsForm.validate()) {
        return
      }

      this.resolve([true, this.name.trim(), this.url.trim()])
      this.resolve = null
      this.show = false
    },

    cancel() {
      this.show = false
    },
  },
}
</script>
