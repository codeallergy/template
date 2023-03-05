<template>
  <section class="section">
    <div class="container">
      <div class="columns">
        <div class="column is-4 is-offset-4">

          <h2 v-if="error.statusCode === 404" class="title has-text-centered">
            {{ pageNotFound }}
          </h2>
          <h2 v-else>
            {{ otherError }}
          </h2>

          <div class="block">{{  this.error }}</div>

        </div>
      </div>   
   </div> 
  </section>     
</template>

<script>
export default {
  name: 'EmptyLayout',
  layout: 'empty',
  props: {
    error: {
      type: Object,
      default: null,
    },
  },
  data() {
    return {
      pageNotFound: '404 Not Found',
      otherError: 'An error occurred',
    }
  },
  head() {
    const title =
      this.error.statusCode === 404 ? this.pageNotFound : this.otherError
    return {
      title,
    }
  },
}
</script>

<style scoped>
h1 {
  font-size: 20px;
}
</style>
