<template>
  <v-app app>
    <v-navigation-drawer app clipped>
      <!-- -->
    </v-navigation-drawer>

    <v-app-bar app clipped-left>
      <v-toolbar-title v-if="user.name" style="font-weight: 700">Hi, {{ user.name }}!</v-toolbar-title>
      <v-spacer></v-spacer>

      <v-btn @click="copyLink" class="mr-2">
        <v-icon>mdi-content-copy</v-icon>
        <span style="margin-left: 10px;">Invite players</span>
      </v-btn>

      <UserNameDialog :user="user"/>

    </v-app-bar>

    <!-- Sizes your content based upon application components -->
    <v-main>
      <router-view :user="user"/>
    </v-main>

    <v-footer app>
      <!-- -->
    </v-footer>


    <InviteDialog ref="inviteDialog" v-model="inviteDialog"></InviteDialog>
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
    openUserDialog: false,
    invitationOpacity: 0,
    inviteDialog: false,
  }),

  created() {
    this.user = new User();
    axios.interceptors.request.use(req => {
      req.headers.authorization = `Bearer ${this.user.id}`;
      return req;
    });
  },

  methods: {
    async copyLink() {
      await this.$refs.inviteDialog.open(location.href)
      this.invitationOpacity = 1
      setTimeout(() => {
        this.invitationOpacity = 0
      }, 3000)
    }

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
