<template>
  <q-page class="row justify-center items-center">
    <div class="column">
      <div class="row">
        <q-card square>
          <q-card-section>
            <q-form class="q-gutter-md" ref="loginForm">
              <q-input
                square
                filled
                v-model="email"
                type="email"
                label="Email"
                lazy-rules
                :rules="[ val => val && val.length > 0 || 'Email is required']"
              />
              <q-input
                square
                filled
                v-model="password"
                type="password"
                label="Password"
                @keydown.enter.prevent="onSubmit"
                lazy-rules
                :rules="[ val => val && val.length > 0 || 'Password is required']"
              />
            </q-form>
          </q-card-section>
          <q-card-actions class="q-px-md">
            <q-btn
              unelevated
              color="primary"
              size="lg"
              class="full-width"
              label="Login"
              @click="onSubmit"
            />
          </q-card-actions>
          <q-card-section class="text-center q-pa-none">
            <p class="text-grey-6">Not reigistered? Created an Account</p>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script>
export default {
  name: 'Login',
  data() {
    return {
      email: '',
      password: ''
    };
  },
  mounted() {
    // When the component is mounted, call to /me endpoint and check if already
    // authenticated
    this.$axios
      .get('http://localhost:8081/api/v1/me')
      .then(response => {
        localStorage.setItem('auth', true);
        this.$store.commit('user/set', response.data);
        this.$router.push('/app');
        this.$q.notify({
          color: 'green-5',
          textColor: 'white',
          icon: 'mdi-check',
          message: 'Already logged in',
          position: 'bottom'
        });
      })
      .catch(() => {});
  },
  methods: {
    onSubmit() {
      this.$refs.loginForm.validate().then(success => {
        if (success) {
          this.$axios
            .post('http://localhost:8081/api/v1/login', {
              email: this.email,
              password: this.password
            })
            .then(response => {
              localStorage.setItem('auth', true);
              this.$router.push('/app');
              this.$q.notify({
                color: 'green-5',
                textColor: 'white',
                icon: 'mdi-check',
                message: 'Successfully logged in',
                position: 'bottom'
              });
            })
            .catch(error => {
              if (error.response) {
                let msg = '';
                switch (error.response.status) {
                  case 401:
                    msg = 'Invalid credentials, wrong email or password';
                    break;
                  default:
                    msg = 'An error occured';
                }
                this.$q.notify({
                  color: 'red-5',
                  textColor: 'white',
                  icon: 'mdi-alert-circle',
                  message: 'Error',
                  caption: msg,
                  position: 'bottom-right'
                });
              }
            });
        }
      });
    }
  }
};
</script>

<style>
.q-card {
  width: 360px;
}
</style>
