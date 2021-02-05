<template>
  <section>
    <v-row no-gutters justify="end" class="pt-3">
      <v-subheader v-text="'active users'" class="font-weight-bold"/>
    </v-row>

    <v-divider/>

    <v-list class="py-0">
      <v-scroll-x-reverse-transition group>
        <v-list-item two-line
                     v-for="item in players"
                     :key="item.Name"
                     @click="() => true"
        >
          <v-list-item-avatar color="primary">
            <v-row no-gutters class="fill-height" align="center" justify="center">
              <span class="font-weight-bold white--text" v-text="getFirstChars(item.Name)"/>
            </v-row>
          </v-list-item-avatar>
          <v-list-item-content>
            <v-list-item-title v-text="item.Name"/>
            <v-list-item-subtitle v-text="state.canReveal() || state.canRestart() ? 'initiator' : 'player'"/>
          </v-list-item-content>

          <v-list-item-action>
            <v-icon v-show="item.VotedCard === '*'"
                    v-text="'mdi-vote-outline'"
                    color="green"
            />
          </v-list-item-action>

        </v-list-item>
      </v-scroll-x-reverse-transition>
    </v-list>
  </section>
</template>

<script>

import State from "@/models/state";

export default {
  name: "UserList.vue",

  data: () => ({
    state: null,
  }),

  computed:{
    players(){
      return this.state ? this.state.Players : [];
    }
  },

  created() {
    State.getInstance(this.setState)
  },

  methods:{
    setState(state){
      this.state = state;
    },

    getFirstChars(str){
      let firstChar = str.replace(' ','').charAt(0),
          secondChar = '';

      const lastSpace = str.lastIndexOf(' ');

      if (lastSpace > 1){
        secondChar = str.charAt( lastSpace + 1);
      }

      return `${firstChar}${secondChar}`.toUpperCase();
    },
  }

}
</script>

<style scoped>

</style>
