<template>
  <div class="columns">
    <div class="column is-6">

      <Notification v-if="error" :message="error"/>
      <FetchNotification :fetchState="$fetchState" />

      <form method="post" autocomplete="off" @submit.prevent="updateUser">

        <article class="message">
          <div class="message-header">
            <p>Edit User</p>
            <button class="delete" aria-label="delete" @click="cancelOperation"></button>
          </div>
          <div class="message-body">

            <div class="block">
              <strong>Id:</strong> {{ userId }}
            </div>

            <div class="block">
              <strong>Email:</strong> {{ email }}
            </div>

            <div class="block">
                <strong>Full Name:</strong> {{ fullName }}
            </div>

            <div class="block">
                <strong>Registered:</strong> {{new Date(createdAt*1000).toLocaleDateString("en-US")}}
            </div>

            <div class="field">
              <label class="label">Role</label>

              <div class="control">
                <div class="select is-primary">
                  <select v-model="role" required>
                    <option value="USER">USER</option>
                    <option value="ADMIN">ADMIN</option>
                  </select>
                </div>
              </div>
            </div>

            <div class="control">
              <button type="submit" class="button is-dark is-fullwidth">Update</button>
            </div>

          </div>
        </article>

      </form>

    </div>
  </div>
</template>

<script>
import FetchNotification from '~/components/FetchNotification';
import Notification from '~/components/Notification';

export default {

  components: {
      FetchNotification,
      Notification,
  },

  layout: 'admin',
  middleware: 'auth-admin',

  data() {
    return {
      userId: '',
      email: '',
      fullName: '',
      role: '',
      createdAt: 0,
      error: null,
    };
  },

  async fetch() {
    const params = this.$route.query
    this.userId = params.id
    const res = await this.$axios.get('/api/admin/users/' + params.id);
    if (res.status === 200) {
      this.userId = res.data.id
      this.email = res.data.email
      this.fullName = res.data.full_name
      this.role = res.data.role
      this.createdAt = res.data.created_at
    }
  },

  watch: {
    '$route.query': '$fetch'
  },

  methods: {
    async updateUser() {
      try {
        await this.$axios.put('/api/admin/users/' + this.userId, {
          role: this.role,
        });

        this.$router.push('/admin/users');
      } catch (e) {
        this.error = e.response.data.message;
      }
    },
    cancelOperation() {
      this.$router.push('/admin/users');
    },
  },

};

</script>


