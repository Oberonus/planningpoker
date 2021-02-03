<template>
  <v-dialog @keydown.enter.prevent="close(true)"
            max-width="500px"
            v-model="openDialog"
            :persistent="!user.identified"
  >
    <template v-slot:activator="{ on, attrs }">
      <v-btn icon v-on="on" v-bind="attrs">
        <v-icon>mdi-account-edit</v-icon>
      </v-btn>
    </template>
    
    <v-card>
      <v-card-title>
        <span class="headline">Choose your name</span>
      </v-card-title>
      <v-divider/>
      <v-card-text class="pt-4">
        <v-text-field
            ref="input"
            v-model="name"
            label="Your display name"
        />
      </v-card-text>
      <v-divider/>
      <v-card-actions>
        <v-spacer/>
        <v-btn color="primary" text @click="close(true)" :disabled="name===''">OK</v-btn>
        <v-btn v-if="user.identified" color="primary" text @click="close(false)">Cancel</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  props: {
    user: {
      required: true,
    },
  },

  data: () => ({
    name: '',
    dialog: false,
  }),

  mounted() {
    this.onDialogOpen();
  },

  computed: {
    openDialog: {
      get: function() {
        return !this.user.identified ? true : this.dialog;
      },
      set: function(value) {
        this.dialog = value;
        value && this.onDialogOpen();
      }
    },
  },

  methods: {
    onDialogOpen(){
      this.name = this.user.name;
      setTimeout(() => this.$refs.input && this.$refs.input.focus(), 0);
    },

    close(ok){
      if (ok) {
        if (this.user.identified){
          this.user.name = this.name;
        } else {
          this.user.name = this.name;
          this.user.authenticate();
        }
      }
      this.dialog = false;
    }
  }
}
</script>
