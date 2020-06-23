<template>
  <q-layout view="hHh lpR fFf">
    <q-header bordered class="bg-primary text-white">
      <q-toolbar>
        <q-btn flat @click="left = !left" round dense icon="mdi-menu" />
        <q-toolbar-title>VueMonit</q-toolbar-title>
        <q-btn
          flat
          round
          :icon="$q.dark.isActive?'mdi-white-balance-sunny':'mdi-weather-night'"
          @click="$q.dark.toggle()"
        />

        <q-btn flat round icon="mdi-refresh" @click="refresh"></q-btn>

        <q-btn flat round icon="mdi-account">
          <q-menu>
            <div class="row no-wrap q-pa-md">
              <div class="column items-center">
                <div class="text-subtitle1 q-mt-md q-mb-xs">{{ user.email }}</div>
                <q-btn color="primary" label="Logout" push size="sm" @click="logout" />
              </div>
            </div>
          </q-menu>
        </q-btn>
      </q-toolbar>
    </q-header>

    <q-drawer
      show-if-above
      v-model="left"
      side="left"
      bordered
      :mini="miniState"
      @mouseover="miniState = false"
      @mouseout="miniState = true"
    >
      <q-list>
        <q-item clickable to="/app" v-ripple>
          <q-item-section avatar>
            <q-icon name="mdi-view-dashboard" />
          </q-item-section>
          <q-item-section>Dashboard</q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <!-- <router-view /> -->
      <q-page class="justify-center items-center">
        <div v-if="loading" class="q-pa-md flex flex-center">
          <q-circular-progress indeterminate size="50px" color="primary" class="q-ma-md" />
        </div>
        <div v-else class="row" v-for="service in services" :key="service.id">
          <q-card flat bordered class="fit q-ma-xs status">
            <div class="fit row inline wrap justify-start items-start content-start">
              <q-item class="col-auto">
                <q-item-section avatar>
                  <q-avatar>
                    <q-icon name="mdi-laptop"></q-icon>
                  </q-avatar>
                </q-item-section>

                <q-item-section>
                  <q-item-label>{{ service.name }}</q-item-label>
                  <q-item-label caption>{{ service.description }}</q-item-label>
                  <q-item-label caption>{{ service.id }}</q-item-label>
                </q-item-section>
              </q-item>
              <q-list class="col-auto">
                <q-item>
                  <q-item-section>
                    <q-item-label>DNS</q-item-label>
                    <q-item-label
                      caption
                      lines="2"
                    >{{ (service.healthcheck.dns / 1000000).toFixed(2) }}ms</q-item-label>
                  </q-item-section>
                </q-item>
                <q-item>
                  <q-item-section>
                    <q-item-label>Handshake</q-item-label>
                    <q-item-label
                      caption
                      lines="2"
                    >{{ (service.healthcheck.handshake / 1000000).toFixed(2) }}ms</q-item-label>
                  </q-item-section>
                </q-item>

                <q-item>
                  <q-item-section>
                    <q-item-label>Connect</q-item-label>
                    <q-item-label
                      caption
                      lines="2"
                    >{{ (service.healthcheck.connect / 1000000).toFixed(2) }}ms</q-item-label>
                  </q-item-section>
                </q-item>
              </q-list>
              <q-list class="col-auto">
                <q-item>
                  <q-item-section>
                    <q-item-label>Total</q-item-label>
                    <q-item-label
                      caption
                      lines="2"
                    >{{ (service.healthcheck.total / 1000000).toFixed(2) }}ms</q-item-label>
                  </q-item-section>
                </q-item>

                <q-item>
                  <q-item-section>
                    <q-item-label>Server</q-item-label>
                    <q-item-label
                      caption
                      lines="2"
                    >{{ (service.healthcheck.server / 1000000).toFixed(2) }}ms</q-item-label>
                  </q-item-section>
                </q-item>
                <q-item>
                  <q-item-section>
                    <q-item-label>Status</q-item-label>
                    <q-item-label caption lines="2">{{ service.healthcheck.status }}</q-item-label>
                  </q-item-section>
                </q-item>
              </q-list>
              <q-list class="col-auto">
                <q-item>
                  <q-item-section>
                    <q-item-label>Metrics Count</q-item-label>
                    <q-item-label caption lines="2">{{ service.count }}</q-item-label>
                  </q-item-section>
                </q-item>

                <q-item>
                  <q-item-section>
                    <q-item-label>Metrics Storage Usage</q-item-label>
                    <q-item-label
                      caption
                      lines="2"
                      class="text-green"
                    >{{ (service.count / 100000).toFixed(2) }}%</q-item-label>
                  </q-item-section>
                </q-item>
              </q-list>
              <q-list>
                <q-item>
                  <q-item-section>
                    <q-item-label>Metrics History</q-item-label>
                    <q-item-label caption lines="2">{{ humanize(service.count*10*1000) }}</q-item-label>
                  </q-item-section>
                </q-item>

                <q-item>
                  <q-item-section>
                    <q-item-label>Remaining History</q-item-label>
                    <q-item-label caption lines="2">{{ humanize((100000 - service.count)*10*1000) }}</q-item-label>
                  </q-item-section>
                </q-item>
              </q-list>
              <div class="col">
                <LineChart :chart-data="service.chartData" :options="chartOpts" :height="100"></LineChart>
              </div>
            </div>
          </q-card>
        </div>
      </q-page>

      <q-dialog v-model="prompt" persistent>
        <q-card style="min-width: 350px">
          <q-card-section>
            <div class="text-h6">Add a Service</div>
          </q-card-section>

          <q-card-section class="q-pt-none">
            <q-input dense autofocus @keyup.enter="prompt = false" />
          </q-card-section>

          <q-card-actions align="right" class="text-primary">
            <q-btn flat label="Cancel" v-close-popup />
            <q-btn flat label="Add Service" v-close-popup />
          </q-card-actions>
        </q-card>
      </q-dialog>
      <q-page-sticky position="bottom-right" :offset="[18, 18]">
        <q-fab icon="mdi-plus" direction="up" color="primary">
          <q-fab-action @click="prompt = true" color="primary" icon="mdi-layers-plus" />
        </q-fab>
      </q-page-sticky>
    </q-page-container>
  </q-layout>
</template>

<script>
import humanizeDuration from 'humanize-duration';
import LineChart from '../components/LineChart.vue';

export default {
  components: {
    LineChart
  },
  data() {
    return {
      loading: true,
      prompt: false,
      miniState: true,
      left: false,
      user: {
        email: ''
      },
      services: [],
      chartOpts: {
        title: {
          display: true,
          text: 'Response Times'
        },
        legend: {
          display: true
        },
        elements: {
          point: {
            radius: 0
          }
        },
        layout: {
          padding: {
            left: 10,
            right: 0,
            top: 0,
            bottom: 10
          }
        },
        tooltips: {
          mode: 'index',
          intersect: false,
          callbacks: {
            label: function(t, d) {
              var label = d.datasets[t.datasetIndex].label || '';
              return label + ': ' + (t.yLabel / 1000000).toFixed(2) + 'ms';
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
              gridLines: {
                display: false
              },
              ticks: {
                autoSkip: true,
                maxTicksLimit: 5,
                beginAtZero: true,
                callback: function(label) {
                  return label / 1000000 + 'ms';
                }
              }
            }
          ]
        }
      }
    };
  },
  methods: {
    humanize: function(input) {
      return humanizeDuration(input, { largest: 2, conjunction: ' and ' });
    },
    historyToChartData: function() {
      let pa = [];
      for (let i = 0; i < this.services.length; i++) {
        const el = this.services[i];
        pa.push(
          this.$axios.get(
            'http://localhost:8081/api/v1/service/' + el.id + '/tr'
          )
        );
      }
      Promise.all(pa)
        .then(results => {
          for (let i = 0; i < pa.length; i++) {
            const res = results[i];
            const el = this.services[i];
            let srvData = [];
            let totalData = [];
            let dates = [];
            for (let i = 0; i < res.data.length && i < 100; i++) {
              srvData[i] = res.data[i].server;
              totalData[i] = res.data[i].total;
              dates[i] = res.data[i].at;
            }
            el.chartData = {
              labels: dates.reverse(),
              datasets: [
                {
                  label: 'Server Response',
                  data: srvData.reverse(),
                  backgroundColor: [['rgba(255, 255, 255, 0.0)']],
                  borderColor: ['#2196f3'],
                  borderWidth: 1
                },
                {
                  label: 'Total Response',
                  data: totalData.reverse(),
                  backgroundColor: [['rgba(255, 255, 255, 0.0)']],
                  borderColor: ['#4caf50'],
                  borderWidth: 1
                }
              ]
            };
          }
        })
        .catch(error => {
          console.log(error);
        })
        .finally(() => {
          this.loading = false;
        });
    },
    formatMs: function(val) {
      console.log(val);
      return (val / 1000000).toFixed(0) + ' ms';
    },
    refresh: function() {
      this.$axios
        .get('http://localhost:8081/api/v1/services')
        .then(response => {
          this.services = response.data;
        })
        .catch(error => {
          this.$q.notify({
            color: 'red-5',
            textColor: 'white',
            icon: 'mdi-alert-circle',
            message: 'Error',
            caption: 'Unable to retrieve services',
            position: 'bottom'
          });
        });
    },
    logout: function() {
      this.$axios
        .get('http://localhost:8081/api/v1/logout')
        .then(response => {
          localStorage.setItem('auth', false);
          this.$router.push('/login');
          this.$q.notify({
            color: 'green-5',
            textColor: 'white',
            icon: 'mdi-check',
            message: 'Successfully logged out',
            position: 'bottom-right'
          });
        })
        .catch(error => {
          this.$q.notify({
            color: 'red-5',
            textColor: 'white',
            icon: 'mdi-alert-circle',
            message: 'Error',
            caption: 'Unable to logout',
            position: 'bottom'
          });
        });
    }
  },
  mounted() {
    this.$axios
      .get('http://localhost:8081/api/v1/me')
      .then(response => {
        this.user.email = response.data.email;
        this.$axios
          .get('http://localhost:8081/api/v1/services')
          .then(response => {
            this.services = response.data;
            this.historyToChartData();
          })
          .catch(error => {
            this.$q.notify({
              color: 'red-5',
              textColor: 'white',
              icon: 'mdi-alert-circle',
              message: 'Error',
              caption: 'Unable to retrieve services',
              position: 'bottom'
            });
          });
      })
      .catch(() => {
        localStorage.setItem('auth', false);
        this.$router.push('/login');
      });
  }
};
</script>

<!-- Notice lang="scss" -->
<style lang="scss">
.q-card.status {
  border-left-width: 4px;
  border-left-color: $green;
  border-radius: 4px;
}
</style>
