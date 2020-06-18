<template>
  <v-app id="inspire">
    <v-navigation-drawer v-model="drawer" app clipped expand-on-hover>
      <v-layout justify-space-between column fill-height>
        <v-list>
          <v-list-item link to="/">
            <v-list-item-action>
              <v-icon>mdi-view-dashboard</v-icon>
            </v-list-item-action>
            <v-list-item-content>
              <v-list-item-title>Dashboard</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list>
        <v-list justify-end>
          <v-list-item link to="/settings">
            <v-list-item-action>
              <v-icon>mdi-cog</v-icon>
            </v-list-item-action>
            <v-list-item-content>
              <v-list-item-title>{{ $t('settings') }}</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-layout>
    </v-navigation-drawer>

    <router-view @drawer="drawer=!drawer" />

    <v-footer app>
      <span>&copy; 2020</span>
    </v-footer>
  </v-app>
</template>

<script>
export default {
  props: {
    source: String
  },
  data: () => ({
    drawer: null
  }),
  created() {
    this.$vuetify.theme.dark = true;
  },
  mounted: function() {
    let elHtml = document.getElementsByTagName("html")[0];
    elHtml.style.overflowY = "auto";
  },
  destroyed: function() {
    let elHtml = document.getElementsByTagName("html")[0];
    elHtml.style.overflowY = null;
  }
};
</script>
