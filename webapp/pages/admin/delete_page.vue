<template>
    <div class="columns">
      <div class="column is-4 is-offset-3">
        <h2 class="title has-text-centered">Delete Page?</h2>

        <Notification v-if="error" :message="error"/>

          <div class="block">
            <strong>Name:</strong> {{ name }}
          </div>

          <div class="block">
             <strong>Title:</strong> {{ title }}
          </div>

          <div class="control">
            <button type="submit" class="button is-dark is-fullwidth" @click="deletePage">Delete</button>
          </div>
        </form>
      </div>
   </div>
   </template>

<script>
import Notification from '~/components/Notification';

export default {

    components: {
        Notification,
    },

    layout: 'admin',
    middleware: 'auth-admin',

    data() {
      return {
        name: '',
        title: '',
        error: null,
      };
    },

    created() {
      this.reloadPage(this.$route.query)
      this.$watch(
          () => this.$route.query,
          (toParams, previousParams) => {
          this.reloadPage(toParams)
          })
    },

    methods: {
      reloadPage(params) {
          this.$axios.get('/api/admin/page/' + params.name)
          .then(res => {
          if(res.status === 200){
              this.name = res.data.name
              this.title = res.data.title
          }
          }).catch((e) => {
              this.error = e.response.data.message;
          })
      },
      async deletePage() {
        try {
          await this.$axios.delete('/api/admin/page/' + this.name);
          this.$router.push('/admin/pages');
        } catch (e) {
          this.error = e.response.data.message;
        }
      },
    },

};

</script>


