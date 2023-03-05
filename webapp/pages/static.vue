<template>
  <section class="section">
    <div class="container">
        <div class="block">
          <h2 v-if="title" class="title has-text-centered">
              {{ title }}
          </h2>
          <iframe
            id="preview"
            ref="preview"
            src="/preview_iframe.html"
            width="100%"
            height="500"
            style="background: white"
            frameborder="0"
            scrolling="yes"
          ></iframe>
        </div>
   </div>
  </section>
</template>

<script>
export default {

  middleware: 'guest',

  data() {
    return {
      title: '',
      content: '',
      error: null,
    };
  },

  created() {
      this.reloadPage(this.$route.query)
      this.$watch(
        () => this.$route.query,
        (toParams, previousParams) => {
          this.reloadPage(toParams)
        })
  },

  methods: {
      reloadPage(params) {
        this.$axios.get('/api/page/' + params.page)
        .then(res => {
          if(res.status === 200){
            this.title = res.data.title
            this.content = res.data.content
            this.updateFrame()
          }
        }).catch((error) => {
          console.log(error)
          this.title = 'Page Not Found'
          this.content = 'Oops, requested page is not found. Try again later.'
        })
      },
      updateFrame() {
          const el = this.$refs.preview.contentWindow.document.getElementById('app')
          if (el != null) {
            el.innerHTML = this.content
          } else {
            setTimeout(this.updateFrame, 10)
          }
      },
  },

}
</script>
