<template>
   <div class="container">

       <div class="columns">
         <div class="column">
             <h2 class="title">Users</h2>
         </div>
       </div>

       <Notification v-if="error" :message="error"/>

       <div v-if="items != null && items.length > 0" class="block">

         <table class="table">
           <thead>
             <tr>
               <th><abbr title="Pos">Pos</abbr></th>
               <th><abbr title="Id">Id</abbr></th>
               <th><abbr title="Email">Email</abbr></th>
               <th><abbr title="Name">Full Name</abbr></th>
               <th><abbr title="Role">Role</abbr></th>
               <th><abbr title="Registered">Registered</abbr></th>
               <th><abbr title="Action">Action</abbr></th>
             </tr>
           </thead>
           <tfoot>
             <tr>
               <th><abbr title="Pos">Pos</abbr></th>
               <th><abbr title="Id">Id</abbr></th>
               <th><abbr title="Email">Email</abbr></th>
               <th><abbr title="Name">Full Name</abbr></th>
               <th><abbr title="Role">Role</abbr></th>
               <th><abbr title="Registered">Registered</abbr></th>
               <th><abbr title="Action">Action</abbr></th>
             </tr>
           </tfoot>
           <tbody>
             <tr v-for="item in items" :key="item.position">
               <th>{{item.position}}</th>
               <th>{{item.id}}</th>
               <td>{{item.email}}</td>
               <td>{{item.full_name}}</td>
               <td>{{item.role}}</td>
               <th>{{new Date(item.created_at*1000).toLocaleDateString("en-US")}}</th>
               <td>
                  <nav class="level">
                    <div class="level-left">
                      <nuxt-link :to="{ path: '/admin/edit_user', query: { id: item.id }}" class="level-item" aria-label="edit">
                        <span class="icon is-small">
                          <font-awesome-icon icon="fa-solid fa-edit" />
                        </span>
                      </nuxt-link>
                      <nuxt-link :to="{ path: '/admin/delete_user', query: { id: item.id }}" class="level-item" aria-label="delete">
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
           const res = await this.$axios.post('/api/admin/users', {
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
         this.$axios.post('/api/admin/users', {
             offset: (page-1) * this.itemsPerPage,
             limit: this.itemsPerPage,
         })
         .then(res => {
           this.items = res.data.items
           this.total = res.data.total
           this.current  = page
         })
       },
     },

 };
</script>
