<template>
  <v-app>
    <router-view></router-view>
  </v-app>
</template>

<script>
export default {
  name: 'App',
};
</script>

<style lang="stylus">
html
body
  //overflow hidden
  overflow-y hidden
  background-color #f0f0f0
</style>