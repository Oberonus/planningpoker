<template>
  <v-app app>

    <v-navigation-drawer app clipped>
      <!-- -->
    </v-navigation-drawer>

    <v-app-bar app clipped-left>
      <v-toolbar-title v-show="user.name"
                       style="font-weight: 700"
                       v-text="`Hi, ${user.name}!`"
      />
      <v-spacer />
      <InviteDialog />
      <UserNameDialog :user="user" />
    </v-app-bar>

    <!-- Sizes your content based upon application components -->
    <v-main>
      <router-view :user="user"/>
    </v-main>

    <v-footer app>
      <!-- -->
    </v-footer>

  </v-app>
</template>

<script>

import User from "@/models/user";
import InviteDialog from "@/components/InviteDialog";
import UserNameDialog from "@/components/UserNameDialog";
import axios from "axios";

export default {
  name: 'App',

  components: {
    UserNameDialog,
    InviteDialog
  },

  data: () => ({
    user: null,
  }),

  created() {
    this.user = new User();
    axios.interceptors.request.use(req => {
      req.headers.authorization = `Bearer ${this.user.id}`;
      return req;
    });
  },

  methods: {

  }
};
</script>

<style lang="stylus">
html
body
  //overflow hidden
  overflow-y hidden
  background-color #f0f0f0
</style>
