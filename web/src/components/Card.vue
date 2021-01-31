<template>
  <div class="card">
    <button
        :class="cardClass"
        @click="$emit('click')"
        translate="no">
      {{ value }}
    </button>
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

<style lang="stylus">
.card
  width: 3rem;
  height: 5rem;
  border-radius: .5rem;
  font-size: 19px;

  margin .4rem .4rem
  display inline-block
  vertical-align top
  white-space nowrap
  transition all .1s linear

  outline: 0;
  cursor: pointer;

.card-btn-inactive:hover
  margin-top: -.005rem;
  box-shadow: 0 3px 5px #888888;
  transition all .1s linear

.card-btn-inactive
  border: 2px solid #23D2AA;

.card-btn-active
  color: #fff;
  background: #20C29A;
  border-color: transparent;
  margin-top: -.8rem;

.card-btn-active:hover
  box-shadow: 0 3px 5px #888888;
  transition all .1s linear

</style>