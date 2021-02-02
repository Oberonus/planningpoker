<template>
  <v-hover v-slot:default="{ hover }" open-delay="100">
    <v-card
        height="6rem"
        width="100%"
        class="rounded-lg card"
        :dark="active"
        :class="hover || active ? cardActiveClasses : cardInactiveClasses"
        :color="active ? activeColor : ''"
        @click.prevent="$emit('select')"
    >
      <v-row
          no-gutters
          align="center"
          justify="center"
          class="fill-height"
      >
        <span v-text="value"
              class="card-label"
              :style="`font-size: ${ active ? '24px' : '19px' }`"
        />
      </v-row>
    </v-card>
  </v-hover>
</template>

<script>

//TODO: develop cards panel in MacOs docker style

export default {
  props: {
    active: Boolean,
    value: String,
  },

  computed: {
    cardActiveClasses() {
      return `card-active ${this.active ? 'card-selected' : ''}`;
    }

  },

  data: () => ({
    activeColor: '#20C29A',
    cardInactiveClasses: 'card-inactive',
  }),

}
</script>

<style lang="stylus">
.card
  transition all .1s linear;
  cursor: pointer;
  border: 2px solid #23D2AA !important;
  overflow: hidden;

.card-inactive
  border: 2px solid transparent !important;

.card-active
  box-shadow: 0 3px 5px #888888 !important;
  transform translateY(-0.5rem);

.card-selected
  transform: translateY(-1rem);

.card-label
  transition: all .1s linear;

</style>
