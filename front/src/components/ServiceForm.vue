<template>
  <q-stepper
    v-model="step"
    ref="stepper"
    color="primary"
    animated
    style="width: 700px; max-width: 80vw;"
  >
    <q-step :name="1" title="Settings" icon="mdi-cog" :done="step > 1">
      <q-form class="q-gutter-md">
        <q-input dense outlined v-model="name" label="Name" hint="Name of the service" />
        <q-input
          dense
          outlined
          v-model="description"
          label="Description"
          hint="Description of the service"
        />
        <q-input
          dense
          outlined
          v-model="healthcheck"
          label="Health Check URL"
          hint="URL to GET for health check"
        />
        <q-input dense outlined v-model="url" label="App URL" hint="URL of the service" />
      </q-form>
    </q-step>
    <q-step :name="2" title="Source" icon="mdi-source-repository" :done="step > 2">
      <q-form class="q-gutter-md">
        <q-toggle v-model="vcs.has" label="Retrieve VCS data" />
        <q-select
          dense
          outlined
          v-if="vcs.has"
          v-model="vcs.system"
          :options="vcs.opts"
          label="VCS Type"
          hint="Type of VCS used"
        />
        <q-input
          v-if="vcs.has"
          v-model="vcs.url"
          dense
          outlined
          label="VCS URL"
          hint="URL of the repo"
        />
        <q-toggle v-if="vcs.has" v-model="vcs.hook" label="Provide Webhook" />
      </q-form>
    </q-step>

    <q-step :name="3" title="CI/CD" icon="mdi-wrench-outline" :done="step > 3">
      <q-form class="q-gutter-md">
        <q-toggle v-model="ci.has" label="Retrieve CI data" />
        <q-select
          dense
          outlined
          v-if="ci.has"
          v-model="ci.system"
          :options="ci.opts"
          label="CI Type"
          hint="Type of VCS used"
        />
        <q-input
          v-if="ci.has"
          v-model="ci.url"
          dense
          outlined
          label="CI URL"
          hint="URL of the service's CI"
        />
      </q-form>
    </q-step>
    <q-step :name="4" title="Summary" icon="mdi-file-outline" :done="step > 4">
      <div class="fit row inline wrap justify-start items-start content-start">
        <q-list class="col-auto">
          <q-item>
            <q-item-section>
              <q-item-label>{{ name }}</q-item-label>
              <q-item-label caption lines="2">Name</q-item-label>
            </q-item-section>
          </q-item>
          <q-item>
            <q-item-section>
              <q-item-label>{{ description }}</q-item-label>
              <q-item-label caption lines="2">Description</q-item-label>
            </q-item-section>
          </q-item>
        </q-list>
      </div>
    </q-step>

    <template v-slot:navigation>
      <q-stepper-navigation>
        <q-btn
          @click="$refs.stepper.next()"
          color="primary"
          :label="step === 4 ? 'Finish' : 'Continue'"
        />
        <q-btn
          v-if="step > 1"
          flat
          color="primary"
          @click="$refs.stepper.previous()"
          label="Back"
          class="q-ml-sm"
        />
      </q-stepper-navigation>
    </template>
  </q-stepper>
</template>

<script>
export default {
  name: 'ServiceForm',
  data() {
    return {
      step: 1,
      name: null,
      description: null,
      healthcheck: null,
      url: null,
      vcs: {
        url: null,
        has: true,
        system: null,
        hook: false,
        opts: ['GitHub', 'Gitlab']
      },
      ci: {
        url: null,
        has: true,
        system: null,
        opts: ['Drone', 'Gitlab', 'Zuul']
      }
    };
  }
};
</script>
