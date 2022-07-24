<template>
  <div class="tcard">
    <div v-if="card===''" class="tcard-body-placeholder"></div>
    <div v-else-if="card==='*'" class="tcard-body-back"></div>
    <div v-else class="tcard-body">

      <v-tooltip top>
        <template v-slot:activator="{ on, attrs }">
          <div v-if="confidence==='high'" v-bind="attrs" v-on="on" class="tcard-label">!</div>
        </template>
        <span>Very confident</span>
      </v-tooltip>

      <v-tooltip top>
        <template v-slot:activator="{ on, attrs }">
          <div v-if="confidence==='low'" v-bind="attrs" v-on="on" class="tcard-label">?</div>
        </template>
        <span>Not sure</span>
      </v-tooltip>

      <div v-if="confidence==='normal'" class="tcard-label"></div>

      {{ card }}
    </div>
    <div class="tcard-name">{{ name }}</div>
  </div>
</template>

<script>
export default {
  props: {
    card: String,
    name: String,
    confidence: String,
  },

  data: () => {
    return {}
  },

  methods: {},
}
</script>

<style lang="stylus">
.tcard-label
  position: absolute
  top: 0
  left: 0
  width: 100%
  height: 25px
  border-top-left-radius .8rem
  border-top-right-radius .8rem
  background #13C29A
  color white
  text-align center
  font-size 17px

.tcard
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-direction: column;
  perspective: 10rem;
  perspective-origin: bottom;
  max-width: 100%;
  margin 0.5rem 0.5rem;

.tcard-body-placeholder
  width: 4rem;
  height: 7rem;
  border-radius: .8rem;
  background: #e8e9ea;

.tcard-body-back
  width: 4rem;
  height: 7rem;
  border-radius: .8rem;
  background: linear-gradient(-135deg, rgb(35, 210, 170) 10%, transparent),
      repeating-linear-gradient(45deg, rgba(35, 210, 170, 1) 0%, rgba(35, 150, 100, 0.6) 5%, transparent 5%, transparent 10%),
      repeating-linear-gradient(-45deg, rgba(35, 210, 170, 0.4) 0%, rgba(35, 150, 100, 0.5) 5%, transparent 5%, transparent 10%);
  background-color: rgba(35, 210, 170, 0.25);

.tcard-body
  font-size 19px
  width: 4rem;
  height: 7rem;
  border-radius: .8rem;
  border: 2px solid #23D2AA;
  white-space nowrap
  display: flex;
  justify-content: center;
  align-items: center;
  font-weight bold

.tcard-name
  font-size: 1rem;
  padding-top: .8rem;
  width: auto;

  display: flex;
  justify-content: center;
  align-items: center;
  text-align: center;
  font-weight: 700;
  line-height: 1.2em;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;

</style>