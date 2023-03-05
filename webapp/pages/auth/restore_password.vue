<template>
  <section class="section">
    <div class="container">
      <div class="columns">
        <div class="column is-4 is-offset-4">
          <h2 class="title has-text-centered">Restore!</h2>

          <Notification v-if="error" :message="error"/>

          <form method="post" @submit.prevent="restore">

            <div class="field">
              <label class="label required">Email</label>

              <div class="control">
                <input
                  v-model="email"
                  type="email"
                  class="input"
                  name="email"
                  required
                >
              </div>
            </div>

            <div class="control">
              <button type="submit" class="button is-dark is-fullwidth">Restore</button>
            </div>
          </form>

          <div class="has-text-centered" style="margin-top: 20px">
            Already have an account? <nuxt-link to="/auth/login">Login</nuxt-link>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
  import Notification from '~/components/Notification';

  export default {

    components: {
      Notification,
    },

    middleware: 'guest',

    data() {
      return {
        email: '',
        error: null,
      };
    },

    methods: {
      async restore() {
        try {
          await this.$axios.post('/api/auth/restore', {
            email: this.email,
          });
          this.$router.push ({path: '/auth/reset_password', query: {email: this.email}})
        } catch (e) {
          this.error = e.response.data.message;
        }
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
