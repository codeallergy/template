<template>
    <div class="container">

        <div class="columns">
          <div class="column">
              <h2 class="title">Pages</h2>
          </div>
          <div class="column is-four-fifths">
            <div class="buttons">
              <button class="button is-primary" @click="createPage">Create Page</button>
            </div>
          </div>
        </div>

        <Notification v-if="error" :message="error"/>

        <div v-if="items != null && items.length > 0" class="block">

          <table class="table">
            <thead>
              <tr>
                <th><abbr title="Pos">Pos</abbr></th>
                <th><abbr title="Name">Name</abbr></th>
                <th><abbr title="Title">Title</abbr></th>
                <th><abbr title="Created">Created</abbr></th>
                <th><abbr title="Action">Action</abbr></th>
              </tr>
            </thead>
            <tfoot>
              <tr>
                <th><abbr title="Pos">Pos</abbr></th>
                <th><abbr title="Name">Name</abbr></th>
                <th><abbr title="Title">Title</abbr></th>
                <th><abbr title="Created">Created</abbr></th>
                <th><abbr title="Action">Action</abbr></th>
              </tr>
            </tfoot>
            <tbody>
              <tr v-for="item in items" :key="item.position">
                <th>{{item.position}}</th>
                <td><nuxt-link :to="{ path: '/static', query: { page: item.name }}">{{item.name}}</nuxt-link></td>
                <td>{{item.title}}</td>
                <th>{{new Date(item.created_at*1000).toLocaleDateString("en-US")}}</th>
                <td>
                  <nav class="level">
                    <div class="level-left">
                      <nuxt-link :to="{ path: '/admin/edit_page', query: { name: item.name }}" class="level-item" aria-label="edit">
                        <span class="icon is-small">
                          <font-awesome-icon icon="fa-solid fa-edit" />
                        </span>
                      </nuxt-link>

                      <nuxt-link :to="{ path: '/admin/delete_page', query: { name: item.name }}" class="level-item" aria-label="delete">
                        <span class="icon is-small">
                          <font-awesome-icon icon="fa-solid fa-trash" />
                        </span>
                      </nuxt-link>
                    </div>
                  </nav>
                </td>
              </tr>
            </tbody>
          </table>

          <Pagination
            :current="current"
            :total="total"
            :itemsPerPage="itemsPerPage"
            :onChange="onChange">
          </Pagination>

        </div>
    </div>
</template>

<script>
  import Pagination from '~/components/Pagination';

  export default {

    components: {
        Notification,
        Pagination,
    },

    layout: 'admin',
    middleware: 'auth-admin',

    data() {
      return {
        items: [],
        current: 1,         // Current page
        total: 0,           // Items total count
        itemsPerPage: 10,   // Items per page
        error: null,
      };
    },

    async created() {
        try {
            const res = await this.$axios.post('/api/admin/pages', {
              offset: 0,
              limit: this.itemsPerPage,
            });
            if(res.status === 200){
              this.items = res.data.items
              this.total = res.data.total
            }
        } catch (e) {
          this.error = e.response.data.message;
        }
      },

      methods: {
        onChange (page) {
          this.$axios.post('/api/admin/pages', {
              offset: (page-1) * this.itemsPerPage,
              limit: this.itemsPerPage,
          })
          .then(res => {
            this.items = res.data.items
            this.total = res.data.total
            this.current  = page
          })
        },
        createPage() {
          this.$router.push('/admin/create_page');
        },
      },

  };
</script>
