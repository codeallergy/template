<template>
  <div>
    <cookie-consent class="notification is-warning" href="/static?page=cookie_policy" />
    <NavigationBar :fixedTop="true"/>
    <nuxt/>
    <Footer/>
  </div>
</template>

<script>
import CookieConsent from 'vue-cookieconsent-component';
import NavigationBar from '~/components/NavigationBar';
import Footer from '~/components/Footer';

export default {
  components: {
    CookieConsent,
    NavigationBar,
    Footer,
  },
  head: {
    bodyAttrs: {
      class: 'has-navbar-fixed-top'
    }
  }
};
</script>
