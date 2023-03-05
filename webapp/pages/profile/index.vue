<template>
  <section class="section">
    <div class="container">
      <h2 class="title">My Profile</h2>

      <div class="content">
        <p>
          <strong>User Id:</strong>
          {{ loggedInUser.user_id }}
        </p>
        <p>
          <strong>First Name:</strong>
          {{ loggedInUser.first_name }}
        </p>
        <p>
          <strong>Middle Name:</strong>
          {{ loggedInUser.middle_name }}
        </p>
        <p>
          <strong>Last Name:</strong>
          {{ loggedInUser.last_name }}
        </p>
        <p>
          <strong>Email:</strong>
          {{ loggedInUser.email }}
        </p>
        <p>
          <strong>Since:</strong>
          {{ loggedInUser.since }} year
        </p>
        <p>
          <strong>Role:</strong>
          {{ loggedInUser.role }}
        </p>
      </div>

      <div class="column is-four-fifths">
        <div class="buttons">
          <button class="button is-primary" @click="refreshTokens">Refresh Tokens</button>
        </div>
      </div>

      <div class="column is-four-fifths">
        <div class="buttons">
          <button class="button is-danger" @click="deleteUser">Delete All</button>
        </div>
      </div>

    </div>
  </section>
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  middleware: 'auth',

  computed: {
    ...mapGetters(['loggedInUser']),
  },

  methods: {
    refreshTokens() {
      console.log(this.$auth)
      this.$auth.refreshTokens()
    },
    deleteUser() {
      this.$router.push('/profile/delete_user');
    },
  },

};
</script>
