<template>
  <div class="columns">
    <div class="column is-5 is-offset-0">
      <h2 class="title has-text-centered">Create Page</h2>

      <Notification v-if="error" :message="error"/>

      <form method="post" @submit.prevent="createPage">

        <div class="field">
          <label class="label required">/static?page={name}</label>

          <div class="control">
            <input
              v-model="name"
              type="text"
              class="input"
              name="page_name"
              required
            />
          </div>
        </div>

        <div class="field">
          <label class="label">Title</label>

          <div class="control">
            <input
              v-model="title"
              type="text"
              class="input"
              name="title"
            />
          </div>
        </div>

        <div class="field">
          <label class="label required">Content</label>

          <div class="control">

            <div class="select is-primary">
              <select v-model="contentType" required>
                <option value="MARKDOWN">Markdown</option>
                <option value="HTML">HTML</option>
              </select>
            </div>
           </div>

          <div class="control" style="margin-top: 5px;">
            <textarea
              v-model="content"
              type="textarea"
              class="textarea is-primary"
              placeholder="Content"
              name="content"
              rows="10"
              required
              @input="updateFrame"
            />
          </div>

        </div>

        <div class="control">
          <button type="submit" class="button is-dark is-fullwidth">Create</button>
        </div>
      </form>
    </div>
    <div class="column">
        <div class="block">
          <h2 class="title has-text-centered">
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
 </div>
 </template>

<script>
import { marked } from 'marked';
import Notification from '~/components/Notification';

export default {

    components: {
        Notification,
    },

    layout: 'admin',
    middleware: 'auth-admin',

    data() {
      return {
        name: '',
        title: '',
        content: '',
        contentType: 'MARKDOWN',
        error: null,
      };
    },

    methods: {
      async createPage() {
        try {
          await this.$axios.post('/api/admin/page', {
            name: this.name,
            title: this.title,
            content: this.content,
            content_type: this.contentType,
          });
          this.$router.push('/admin/pages');
        } catch (e) {
          this.error = e.response.data.message;
        }
      },
      updateFrame() {
         let htmlContent = this.content
         if (this.contentType === 'MARKDOWN') {
            htmlContent = marked.parse(htmlContent)
         }
         this.$refs.preview.contentWindow.document.getElementById('app').innerHTML = htmlContent
      },
    },

};

</script>

<style scoped>
 .required:after {
   content:" *";
   color: red;
 }

</style>
