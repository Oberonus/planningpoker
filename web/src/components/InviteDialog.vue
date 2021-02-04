<template>
  <v-dialog v-model="dialog" max-width="500px">

    <template v-slot:activator="{ on, attrs }">
      <v-btn v-on="on" v-bind="attrs" class="mr-1">
        <v-icon left v-text="'mdi-content-copy'"/>
        <span v-text="'Invite players'"/>
      </v-btn>
    </template>

    <v-card>
      <v-card-title>
        <span class="headline">Invite your teammates</span>
      </v-card-title>

      <v-divider></v-divider>

      <v-card-text class="mt-5">
        <v-text-field ref="textToCopy"
                      readonly
                      :value="url"
                      @click.prevent="open"
                      label="URL to copy"/>
      </v-card-text>

      <v-divider/>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary" text @click="copy">
          <span v-text="'Copy invitation link'"/>
        </v-btn>
      </v-card-actions>
    </v-card>

  </v-dialog>
</template>

<script>
export default {
  data: () => ({
    dialog: false,
    url: '',
  }),

  watch: {
    dialog: function(value) {
      value && this.open();
    },
  },

  methods: {
    open() {
      this.url = location.href;
      setTimeout(() => {
        let textToCopy = this.$refs.textToCopy.$el.querySelector('input');
        textToCopy.select();
      }, 0)
    },

    copy() {
      let textToCopy = this.$refs.textToCopy.$el.querySelector('input')
      textToCopy.select()
      document.execCommand("copy");
      this.dialog = false;
    },

  },
}
</script>
