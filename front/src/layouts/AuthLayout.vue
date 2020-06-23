<template>
  <q-layout view="hHh lpR fFf">
    <q-header bordered class="bg-primary text-white">
      <q-toolbar>
        <q-toolbar-title>VueMonit</q-toolbar-title>
        <q-btn
          flat
          round
          dense
          :icon="$q.dark.isActive?'mdi-white-balance-sunny':'mdi-weather-night'"
          @click="$q.dark.toggle()"
        />
      </q-toolbar>
    </q-header>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script>
</script>

<style>
</style>
