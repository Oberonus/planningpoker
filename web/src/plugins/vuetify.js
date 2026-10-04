import Vue from 'vue';
import Vuetify from 'vuetify/lib/framework';

Vue.use(Vuetify);

export default new Vuetify({
  theme: {
    themes: {
      light: {
        primary: '#16a382',
        secondary: '#646464',
        error: '#b43d3d',
        success: '#126b52',
      },
    },
  },
});
