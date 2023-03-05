<template>
  <div>
    <NavigationBar/>

    <section class="section">
      <div class="container">
        <div class="columns">
           <div class="column is-one-quarter">
              <aside class="menu">
                <p class="menu-label">
                  General
                </p>
                <ul class="menu-list">
                  <li><nuxt-link to="/admin/">Dashboard</nuxt-link></li>
                  <li><nuxt-link to="/admin/users">Users</nuxt-link></li>
                </ul>
                <p class="menu-label">
                  Administration
                </p>
                <ul class="menu-list">
                  <li><nuxt-link to="/admin/pages">Pages</nuxt-link></li>
                </ul>
                <p class="menu-label">
                  Statistics
                </p>
                <ul class="menu-list">
                  <li><nuxt-link to="/admin/traffic">Traffic</nuxt-link></li>
                </ul>
              </aside>
           </div>
           <div class="column">
              <nuxt/>
           </div>
        </div>
      </div>
    </section>

    <Footer/>
  </div>
</template>

<script>
  import NavigationBar from '~/components/NavigationBar';
  import Footer from '~/components/Footer';

  export default {
    components: {
      NavigationBar,
      Footer,
    },
  };
</script>
