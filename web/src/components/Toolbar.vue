<template>
  <div class="v-toolbar">

  </div>
</template>

<script>
export default {
  props: {
    active: Boolean,
    value: String,
  },

  data: () => {
    return {}
  },

  computed: {
    cardClass() {
      return this.active ? "card card-btn-active" : "card card-btn-inactive"
    }
  },

  methods: {},
}
</script>
