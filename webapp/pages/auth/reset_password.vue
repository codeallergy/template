<template>
  <section class="section">
    <div class="container">
      <div class="columns">
        <div class="column is-4 is-offset-4">
          <h2 class="title has-text-centered">New Password!</h2>

          <Notification v-if="error" :message="error"/>

          <form method="post" autocomplete="off" @submit.prevent="restore">

            <div class="field">
              <label class="label required">Recovery Code</label>

              <div class="control">
                <input
                  v-model="passcode"
                  type="text"
                  class="input"
                  name="passcode"
                  required
                >
              </div>
            </div>

            <div class="field">
              <label class="label required">New Password</label>

              <div class="control">
                <VuePassword
                    v-model="password"
                    :strength="strength"
                    type="password"
                    required
                    @input="updateStrength"
                />
              </div>
            </div>

            <div class="control">
              <button :disabled="strength === 0" type="submit" class="button is-dark is-fullwidth">Reset</button>
            </div>
          </form>

          <div class="has-text-centered" style="margin-top: 20px">
            Already got an account? <nuxt-link to="/auth/login">Login</nuxt-link>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
  import VuePassword from 'vue-password';
  import Notification from '~/components/Notification';

  export default {

    components: {
      VuePassword,
      Notification,
    },

    middleware: 'guest',

    data() {
      return {
        passcode: '',
        password: '',
        strength: 0,

        email: '',
        error: null,
      };
    },

    fetch() {
        if ('email' in this.$route.query) {
          this.email = this.$route.query.email;
          this.error = 'Recovery code was sent to email ' + this.email;
        }
    },

    methods: {
      async restore() {
        try {
          await this.$axios.post('/api/auth/reset', {
            email: this.email,
            code: this.passcode,
            password: this.password,
          });

          await this.$auth.loginWith('local', {
            data: {
              email: this.email,
              password: this.password,
            },
          });

          this.$router.push('/');
        } catch (e) {
          this.error = e.response.data.message;
        }
      },
      unrepeated(str) {
         return [...new Set(str)].join('');
      },
      updateStrength(password) {
        const p = this.unrepeated(password);
        let score = Math.floor(p.length / 3);
        // The password has at least one uppercase letter
        if (p.match(/[A-Z]/g)) {
          score = score + 1;
        }
        // The password has at least one lowercase letter
        if (p.match(/[a-z]/g)) {
          score = score + 1;
        }
        // The password has at least one digit
        if (p.match(/[0-9]/g)) {
          score = score + 1;
        }
        // The password has at least one special character
        if (p.match(/[^A-Za-z0-9]/g)) {
          score = score + 1;
        }
        score = Math.floor(score / 2);
        if (score > 4) {
           score = 4;
        }
        this.strength = score;
      },
    },
  };
</script>

<style>
  .required:after {
    content:" *";
    color: red;
  }
  .red {
    color: red;
  }
</style>
