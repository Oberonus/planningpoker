<template>
  <v-dialog v-model="show" max-width="500px">
    <v-card>
      <v-card-title>
        <span class="headline">Invite your teammates</span>
      </v-card-title>
      <v-divider></v-divider>
      <v-card-text style="margin-top: 20px;">
        <v-row>
          <v-col cols="12">
            <v-text-field ref="textToCopy" v-model="url" label="URL to copy" requred></v-text-field>
          </v-col>
        </v-row>
      </v-card-text>
      <v-divider></v-divider>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary" text @click="copy">Copy invitation link</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  data: () => {
    return {
      show: false,
      url: '',
      resolve: null,
    }
  },

  methods: {
    async open(url) {
      this.url = url
      this.show = true

      this.$nextTick(() => {
        let textToCopy = this.$refs.textToCopy.$el.querySelector('input')
        textToCopy.select()
      })

      return new Promise((resolve) => {
        this.resolve = resolve
      })
    },

    copy() {
      let textToCopy = this.$refs.textToCopy.$el.querySelector('input')
      textToCopy.select()
      document.execCommand("copy");
      this.show = false
      this.resolve()
    },
  },
}
</script>