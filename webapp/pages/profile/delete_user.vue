<template>
  <div class="columns section">
    <div class="column is-4 is-offset-3">
      <h2 class="title has-text-centered">Delete Confirmation</h2>

      <Notification v-if="error" :message="error"/>

      <form method="post" autocomplete="off" @submit.prevent="deleteUser">

        <article class="message">
          <div class="message-header">
            <p>Do you want to delete your user account?</p>
            <button class="delete" aria-label="delete" @click="cancelOperation"></button>
          </div>
          <div class="message-body">
            <p>Deletion of user account action is non-revertable and after clicking on button <strong>Delete All</strong>
            your account will revoke all memberships and permission and all data would be removed from our databases.</p>
            <br/>
            <p>Your account is registered on user <strong>{{loggedInUser.firstName}} {{ loggedInUser.lastName }}</strong> with user id <strong>{{loggedInUser.userId}}</strong>.</p>
            <br/>
            <div class="control">
              <button type="submit" class="button is-danger is-fullwidth">Delete All</button>
            </div>

          </div>
        </article>

      </form>

    </div>
 </div>
 </template>

<script>
import { mapGetters } from 'vuex';
import Notification from '~/components/Notification';

export default {

  components: {
    Notification,
  },

  layout: 'default',
  middleware: 'auth',

  data() {
    return {
      error: null,
    };
  },

  computed: {
    ...mapGetters(['loggedInUser']),
  },

  methods: {
    async deleteUser() {
      try {
        await this.$axios.delete('/api/user/' + this.loggedInUser.user_id);
        await this.$auth.logout();
        this.$router.push('/');
      } catch (e) {
        this.error = e.response.data.message;
      }
    },
    cancelOperation() {
      this.$router.push('/profile');
    },
  }
};

</script>
