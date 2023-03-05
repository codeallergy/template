<template>
  <section class="section">

    <div v-if="isAuthenticated" class="container">

      <section class="hero is-fullheight-with-navbar">
        <div class="hero-body">
          <div class="container has-text-centered">
            <p class="title">
              PreCook Template is changing how the world create web applications
            </p>
            <p class="subtitle">
              Authenticated
            </p>
          </div>
        </div>
      </section>

    </div>

    <div v-else class="container">

      <section class="hero is-fullheight-with-navbar">
        <div class="hero-body">
          <div class="container has-text-centered">
            <p class="title">
              PreCook Template is changing how the world create web applications
            </p>
            <p class="subtitle">
              Not Authenticated
            </p>
          </div>
        </div>
      </section>

    </div>
  </section>
</template>


<script>
import { mapGetters } from 'vuex';

export default {
  computed: {
    ...mapGetters(['isAuthenticated', 'loggedInUser']),
  },
};
</script>
