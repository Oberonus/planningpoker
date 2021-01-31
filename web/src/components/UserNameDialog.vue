<template>
  <v-dialog v-model="show" max-width="500px" persistent>
    <v-card>
      <v-card-title>
        <span class="headline">Choose your name</span>
      </v-card-title>
      <v-divider></v-divider>
      <v-card-text style="margin-top: 20px;">
        <v-row>
          <v-col cols="12">
            <v-text-field v-model="name" label="Your display name" requred></v-text-field>
          </v-col>
        </v-row>
      </v-card-text>
      <v-divider></v-divider>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary" text @click="save" :disabled="name===''">OK</v-btn>
        <v-btn v-if="mode===modeModify" color="primary" text @click="cancel">Cancel</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  data: () => {
    return {
      modeCreate: 'create',
      modeModify: 'modify',
      mode: null,
      show: false,
      name: '',
      resolve: null,
    }
  },

  methods: {
    async open(name) {
      this.mode = this.modeCreate
      this.name = name
      this.show = true

      return new Promise((resolve) => {
        this.resolve = resolve
      })
    },

    async openModify(name) {
      this.mode = this.modeModify
      this.name = name
      this.show = true

      return new Promise((resolve) => {
        this.resolve = resolve
      })
    },

    save() {
      this.show = false
      this.resolve(this.name)
    },

    cancel() {
      this.show = false
      this.resolve("")
    }
  },
}
</script>