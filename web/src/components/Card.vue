<template>
  <div class="card">
    <div :class="cardClass" @click="$emit('click')" translate="no">
      <div v-if="active && confidence==='high'" class="card-btn-active-label">!</div>
      <div v-if="active && confidence==='low'" class="card-btn-active-label">?</div>
      <div v-if="active && confidence==='normal'" class="card-btn-active-label"></div>
      <div class="card-content">{{ value }}</div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    active: Boolean,
    value: String,
    confidence: String,
  },

  data: () => {
    return {}
  },

  computed: {
    cardClass() {
      if (!this.active) {
        return "card card-btn-inactive"
      }
      return "card card-btn-active"
    },
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
  position relative
  vertical-align top
  white-space nowrap
  transition all .1s linear

  outline: 0;
  cursor: pointer;

.card-content
  display: flex;
  align-items: center;
  justify-content center;
  flex-direction: column;
  width: 2.7rem;
  height: 5rem;

.card-btn-inactive:hover
  margin-top: -.005rem;
  box-shadow: 0 3px 5px #888888;
  transition all .1s linear

.card-btn-inactive
  border: 2px solid #A0A0A0;

.card-btn-active
  border: 2px solid #13C29A;
  margin-top: -.5rem;

.card-btn-active:hover
  box-shadow: 0 3px 5px #888888;
  transition all .1s linear

.card-btn-active-label
  position absolute
  top 0
  left 0
  width 100%
  height 20px
  border-top-left-radius .3rem
  border-top-right-radius .3rem
  background #13C29A
  color white
  font-size 13px
  text-align center
  font-weight bold

</style>