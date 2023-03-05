<template>
  <footer class="footer">
    <div class="content has-text-centered">
      <p>
        <strong>Precook Template</strong> by <a href="https://schwid.com">Zander Schwid & Co. LLC</a>. All rights reserved.
      </p>
    </div>

      <div class="bd-footer-links">
        <div class="columns">
            <div class="column is-3 is-offset-4">
              <p class="bd-footer-link-title"><nuxt-link class="navbar-item" to="/">Home</nuxt-link></p>
              <p class="bd-footer-link"><nuxt-link class="navbar-item" :to="{ path: '/static', query: { page: 'about' }}">About</nuxt-link></p>
              <p class="bd-footer-link"><nuxt-link class="navbar-item" :to="{ path: '/static', query: { page: 'terms_of_use' }}">Terms Of Use</nuxt-link></p>
              <p class="bd-footer-link"><nuxt-link class="navbar-item" :to="{ path: '/static', query: { page: 'privacy' }}">Privacy</nuxt-link></p>
              <p class="bd-footer-link"><nuxt-link class="navbar-item" :to="{ path: '/static', query: { page: 'cookie_policy' }}">Cookie Policy</nuxt-link></p>
            </div>
            <div class="column is-3" style="text-align: left;">
              <p class="bd-footer-link-title"><span class="navbar-item">Share on social media</span></p>
              <p class="bd-footer-link"><a class="navbar-item" href="https://twitter.com/share">Tweet</a></p>
            </div>
        </div>
      </div>


  </footer>
</template>

<style>

.bd-footer-links {
  padding-bottom: 6rem;
}

.bd-footer-link-title:not(:first-child) {
    margin-top: 1.5em;
}
.bd-footer-link-title {
    color: #363636;
    font-size: 1.25rem;
    font-weight: 600;
}

</style>

