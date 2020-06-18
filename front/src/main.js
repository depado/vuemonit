import Vue from "vue";
import App from "./App.vue";
import vuetify from "./plugins/vuetify";
import router from "./router";
import i18n from "./i18n";

Vue.config.productionTip = false;
Vue.mixin({
  methods: {
    apiURL: function() {
      return process.env.VUE_APP_API;
    },
  },
});

new Vue({
  vuetify,
  router,
  i18n,
  render: (h) => h(App),
}).$mount("#app");
