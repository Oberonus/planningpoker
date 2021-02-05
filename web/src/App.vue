<template>
  <v-app app>

    <v-navigation-drawer v-model="drawer" app clipped>
      <user-list/>
    </v-navigation-drawer>

    <v-app-bar app clipped-left>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"/>
      <v-toolbar-title v-show="user.name"
          style="font-weight: 700"
          v-text="`Hi, ${user.name}!`"
      />
      <v-spacer />
      <InviteDialog @notify="snackbarNotify"/>
      <UserNameDialog :user="user" @notify="snackbarNotify"/>

      <v-snackbar top right absolute
                  :timeout="2000"
                  v-model="snackbar"
                  @click="snackbar=false"
      >
        <v-row no-gutters justify="space-between" align="center">
          <span v-text="snackbarMessage"/>
          <v-btn icon x-small @click="snackbar=false">
            <v-icon v-text="'mdi-close'"/>
          </v-btn>
        </v-row>
      </v-snackbar>

    </v-app-bar>

    <!-- Sizes your content based upon application components -->
    <v-main>
      <v-fade-transition>
        <router-view :user="user"/>
      </v-fade-transition>
    </v-main>

    <v-footer app>
      <!-- -->
    </v-footer>

  </v-app>
</template>

<script>

import User from "@/models/user";
import UserList from "@/components/UserList";
import InviteDialog from "@/components/InviteDialog";
import UserNameDialog from "@/components/UserNameDialog";
import axios from "axios";

export default {
  name: 'App',

  components: {
    UserNameDialog,
    InviteDialog,
    UserList
  },

  data: () => ({
    user: null,
    drawer: true,
    snackbar: false,
    snackbarMessage: '',
  }),

  created() {
    this.user = new User();
    axios.interceptors.request.use(req => {
      req.headers.authorization = `Bearer ${this.user.id}`;
      return req;
    });
  },

  methods: {
    snackbarNotify(event){
      if (event) {
        this.snackbarMessage = event;
        this.snackbar = true;
      }
    }
  }
};
</script>

<style lang="stylus">
html
body
  overflow-y hidden
  background-color #f0f0f0
</style>
