<template>
  <div>
    <v-app-bar app clipped-left>
      <v-app-bar-nav-icon @click.stop="$emit('drawer')"></v-app-bar-nav-icon>
      <v-toolbar-title>Dashboard</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-tooltip left>
        <template v-slot:activator="{ on }">
          <v-btn v-on="on" @click="refresh" icon>
            <v-icon>mdi-refresh</v-icon>
          </v-btn>
        </template>
        <span>{{ $t('refresh') }}</span>
      </v-tooltip>
    </v-app-bar>
    <v-content>
      <v-container fluid>
        <v-row>
          <v-col cols="6" md="6" v-for="service in services" :key="service.id">
            <v-card>
              <v-card-title>{{ service.name }}</v-card-title>
              <v-card-subtitle>{{ service.healthcheck.url }}</v-card-subtitle>
              <v-card-text>
                <b>DNS:</b>
                {{ (service.healthcheck.dns / 1000000).toFixed(2) }}ms
                <br />
                <b>Handshake:</b>
                {{ (service.healthcheck.handshake / 1000000).toFixed(2) }}ms
                <br />
                <b>Connect:</b>
                {{ (service.healthcheck.connect / 1000000).toFixed(2) }}ms
                <br />
                <b>Total:</b>
                {{ (service.healthcheck.total / 1000000).toFixed(2) }}ms
                <br />
                <b>Server</b>
                {{ (service.healthcheck.server / 1000000).toFixed(2) }}ms
                <br />
              </v-card-text>
              <LineChart :chart-data="service.chartData" :options="chartOpts" :height="100"></LineChart>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-content>
  </div>
</template>

<script>
import axios from "axios";
import LineChart from "../components/LineChart.vue";

export default {
  name: "Dashboard",
  components: {
    LineChart
  },
  data() {
    return {
      services: [],
      loading: true,
      chartOpts: {
        legend: {
          display: false
        },
        tooltips: {
          callbacks: {
            label: function(tooltipItem) {
              return tooltipItem.yLabel;
            }
          }
        },
        scales: {
          xAxes: [
            {
              display: false
            }
          ],
          yAxes: [
            {
              ticks: {
                beginAtZero: true,
                callback: function(label) {
                  return label / 1000000 + "ms";
                }
              }
            }
          ]
        }
      }
    };
  },
  methods: {
    historyToChartData: function() {
      for (let i = 0; i < this.services.length; i++) {
        const el = this.services[i];
        let srvData = [];
        let totalData = [];
        let dates = [];
        for (let i = 0; i < el.history.length; i++) {
          srvData[i] = el.history[i].server;
          totalData[i] = el.history[i].total;
          dates[i] = el.history[i].at;
        }
        el.chartData = {
          labels: dates,
          datasets: [
            {
              label: "Server Response",
              data: srvData,
              backgroundColor: ["rgba(255, 99, 132, 0.2)"],
              borderColor: ["rgba(255, 99, 132, 1)"],
              borderWidth: 1
            },
            {
              label: "Total Response",
              data: totalData,
              borderColor: ["rgba(255, 99, 132, 1)"]
            }
          ]
        };
      }
    },
    formatMs: function(val) {
      console.log(val);
      return (val / 1000000).toFixed(0) + " ms";
    },
    refresh: function() {
      this.loading = true;
      axios
        .get(this.apiURL() + "/fetch")
        .then(response => {
          this.services = response.data;
          this.historyToChartData();
        })
        .catch(error => {
          console.log(error);
        })
        .finally(() => (this.loading = false));
    }
  },
  mounted() {
    this.refresh();
  }
};
</script>
