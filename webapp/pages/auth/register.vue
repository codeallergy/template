<template>
  <section class="section">
    <div class="container">
      <div class="columns">
        <div class="column is-4 is-offset-4">
          <h2 class="title has-text-centered">Register!</h2>

          <Notification v-if="error" :message="error"/>

          <form method="post" @submit.prevent="register">

            <div class="field">
              <label class="label required">First Name</label>

              <div class="control">
                <input
                  v-model="firstName"
                  type="text"
                  class="input"
                  name="first_name"
                  required
                />
              </div>
            </div>

            <div class="field">
              <label class="label">Middle Name</label>

              <div class="control">
                <input
                  v-model="middleName"
                  type="text"
                  class="input"
                  name="middle_name"
                />
              </div>
            </div>

            <div class="field">
              <label class="label required">Last Name</label>

              <div class="control">
                <input
                  v-model="lastName"
                  type="text"
                  class="input"
                  name="last_name"
                  required
                />
              </div>
            </div>

            <div class="field">
              <label class="label required">Email</label>

              <div class="control">
                <input
                  v-model="email"
                  type="email"
                  class="input"
                  name="email"
                  required
                />
              </div>
            </div>

            <div class="field">
              <label class="label required">Password</label>

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

            <br/>

            <div class="container">
              <label class="checkbox">
                <input v-model="agree" type="checkbox" required />
                I agree to the <nuxt-link :to="{ path: '/static', query: { page: 'terms_of_use' }}" target="_blank">Terms Of Use</nuxt-link>
              </label>
            </div>

            <br/>

            <div class="control">
              <button :disabled="strength === 0" type="submit" class="button is-dark is-fullwidth">Register</button>
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
        firstName: '',
        middleName: '',
        lastName: '',
        email: '',
        password: '',
        strength: 0,
        agree: false,
        error: null,
      };
    },

    methods: {
      async register() {
        try {
          await this.$axios.post('/api/auth/register', {
            first_name: this.firstName,
            middle_name: this.middleName,
            last_name: this.lastName,
            email: this.email,
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
</style>
