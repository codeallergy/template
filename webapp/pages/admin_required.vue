<template>
<div class="section">
  <div class="columns">
    <div class="column is-4 is-offset-3">
      <h2 class="title has-text-centered">Oops</h2>

      <br/>
      <div class="notification is-danger">Admin role is required.</div>
      <br/>

    </div>
  </div>
</div>
</template>

<script>

export default {

  layout: 'default',
  middleware: 'auth',

  data() {
    return {
    };
  },


};

</script>
