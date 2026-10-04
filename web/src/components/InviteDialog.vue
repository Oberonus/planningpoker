<template>
  <v-dialog v-model="show" max-width="480px">
    <v-card class="dialog-card" role="dialog" aria-modal="true" aria-labelledby="invite-dialog-title" @keydown.esc.native.stop.prevent="show = false">
      <v-card-title id="invite-dialog-title" role="heading" aria-level="2">Invite your team</v-card-title>
      <v-card-text>
        <p class="dialog-description">Share this link so your teammates can join the game.</p>
        <v-text-field
          ref="textToCopy"
          v-model="url"
          label="Invitation link"
          readonly
          outlined
          dense
          @focus="selectLink"
        ></v-text-field>
        <p v-if="copyFailed" class="muted" role="status">Select the link and copy it with your keyboard or your device's copy menu.</p>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn text @click="show = false">Close</v-btn>
        <v-btn color="primary" depressed @click="copy">Copy link</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  data: () => ({
    show: false,
    url: '',
    resolve: null,
    copyFailed: false,
  }),

  watch: {
    show(visible) {
      if (!visible && this.resolve) {
        this.resolve(false)
        this.resolve = null
      }
    },
  },

  methods: {
    open(url) {
      this.url = url
      this.copyFailed = false
      this.show = true
      this.$nextTick(() => {
        requestAnimationFrame(() => {
          if (this.show) {
            this.selectLink()
          }
        })
      })

      return new Promise(resolve => {
        this.resolve = resolve
      })
    },

    selectLink() {
      const input = this.$refs.textToCopy.$el.querySelector('input')
      input.focus()
      input.select()
    },

    copy() {
      this.selectLink()

      try {
        // Older browsers and non-secure local hosts need the selected-input fallback.
        if (!document.execCommand('copy')) {
          this.copyFailed = true
          return
        }
      } catch (error) {
        this.copyFailed = true
        return
      }

      this.resolve(true)
      this.resolve = null
      this.show = false
    },
  },
}
</script>
