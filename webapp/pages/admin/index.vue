<template>
<div class="tile is-ancestor">
  <div class="tile is-vertical is-8">
    <div class="tile">
      <div class="tile is-parent is-vertical">
        <article class="tile is-child notification is-primary">
          <p class="title">Users</p>
          <p class="subtitle">{{ users }}</p>
        </article>
        <article class="tile is-child notification is-warning">
          <p class="title">Sites</p>
          <p class="subtitle">{{ sites }}</p>
        </article>
      </div>
      <div class="tile is-parent">
        <article class="tile is-child notification is-info">
          <p class="title">Registrations</p>
          <p class="subtitle">New users per day {{ regPerDay }}</p>
        </article>
      </div>
    </div>
    <div class="tile is-parent">
      <article class="tile is-child notification is-danger">
        <p class="title">Traffic</p>
        <p class="subtitle">Hits per minute {{  hitsPerMinute }}</p>
        <div class="content">
          <!-- Content -->
        </div>
      </article>
    </div>
  </div>
  <div class="tile is-parent">
    <article class="tile is-child notification is-success">
      <div class="content">
        <p class="title">Top 10 Users</p>
        <p class="subtitle">Here it is</p>
        <div class="content">
          <!-- Content -->
        </div>
      </div>
    </article>
  </div>
</div>
</template>

<script>

  export default {
    layout: 'admin',
    middleware: 'auth-admin',

    data() {
      return {
        users: 1,
        sites: 1,
        hitsPerMinute: 1,
        regPerDay: 1,
        error: null,
      };
    },

  };
</script>
