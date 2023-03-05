<template>
    <section class="section">
      <div class="container">
        <h2 class="title">Security Log</h2>

          <Notification v-if="error" :message="error"/>

          <div v-if="items != null && items.length > 0" class="block">

            <table class="table">
              <thead>
                <tr>
                  <th><abbr title="Position">Pos</abbr></th>
                  <th><abbr title="Event Name">Event</abbr></th>
                  <th><abbr title="Event Time">Time</abbr></th>
                  <th><abbr title="Remote IP">IP</abbr></th>
                  <th><abbr title="User Agent">User Agent</abbr></th>
                </tr>
              </thead>
              <tfoot>
                <tr>
                  <th><abbr title="Position">Pos</abbr></th>
                  <th><abbr title="Event Name">Event</abbr></th>
                  <th><abbr title="Event Time">Time</abbr></th>
                  <th><abbr title="Remote IP">IP</abbr></th>
                  <th><abbr title="User Agent">User Agent</abbr></th>
                </tr>
              </tfoot>
              <tbody>
                <tr v-for="item in items" :key="item.position">
                  <th>{{item.position}}</th>
                  <td><strong>{{item.event_name}}</strong></td>
                  <th>{{new Date(item.event_time*1000).toLocaleString("en-US")}}</th>
                  <td>{{item.remote_ip}}</td>
                  <td>{{item.user_agent}}</td>
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
    </section>
  </template>

  <script>
    import { mapGetters } from 'vuex';
    import Pagination from '~/components/Pagination';

    export default {

      components: {
        Notification,
        Pagination,
      },

      middleware: 'auth',

      data() {
        return {
          items: [],
          current: 1,        // Current page
          total: 0,          // Items total count
          itemsPerPage: 10,   // Items per page
          error: null,
        };
      },

      computed: {
        ...mapGetters(['loggedInUser']),
      },

      async created() {
        try {
            const res = await this.$axios.post('/api/auth/security_log', {
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
          this.$axios.post('/api/auth/security_log', {
              offset: (page-1) * this.itemsPerPage,
              limit: this.itemsPerPage,
          })
          .then(res => {
            this.items = res.data.items
            this.total = res.data.total
            this.current  = page
          })
        }
      },

    };
  </script>
