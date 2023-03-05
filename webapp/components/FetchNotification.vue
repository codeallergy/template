<template>
  <div>
    <p v-if="showPending && fetchState.pending">Fetching ...</p>
    <p v-if="fetchState.error" class="notification is-danger">
      {{ fetchState.error.response.data.message }}
    </p>
  </div>
</template>

<script>
export default {
  name: 'FetchNotification',
  props: {
    fetchState: Object,
    showPending: Boolean,
  }
};
</script>
