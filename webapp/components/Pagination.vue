<template>
  <div id="pagination">
    <nav class="pagination is-centered">
      <a v-if="current > 1" class="pagination-previous" @click="onChange(current - 1)">&larr; Previous</a>
      <a  v-if="size > 1 && current < size" class="pagination-next" @click="onChange(current + 1)">Next &rarr;</a>
      <ul class="pagination-list">
        <component :is="element.type" v-for="element in elements" :key="element.page" :page="element.page" :current="current" :onChange="onChange"/>
      </ul>
    </nav>
  </div>
</template>

<script>
import PageSquare from './PageSquare.vue'
import EllipseBreak from './EllipseBreak.vue'

export default {
  name: 'Pagination',
  components: { PageSquare, EllipseBreak },
  props: {
    current: {
      type: Number,
      default: 0
    },
    total: {
      type: Number,
      default: 0
    },
    itemsPerPage: {
      type: Number,
      default: 0
    },
    step: {
      type: Number,
      default: 3
    },
    onChange: {
      type: Function
    }
  },
  data () {
    return {
      elements: [],
      size: 0
    }
  },
  mounted () {
    this.paginate()
    // Check for changes in props resulting from asynchronous operations
    Object.keys(this._props).forEach(event => {
      this.$watch(event, (val, oldVal) => {
        this.paginate()
      });
    })
  },
  methods: {
    add (s, f) {
      for (let i = s; i < f; i++) {
        this.elements.push(
          { type: 'page-square', page: i }
        )
      }
    },
    first () {
      // Add first page with separator
      this.elements.push(
        { type: 'page-square', page: 1 },
        { type: 'ellipse-break' }
      )
    },
    last () {
      // Add last page with separator
      this.elements.push(
        { type: 'ellipse-break' },
        { type: 'page-square', page: this.size },
      )
    },
    paginate () {
      this.elements = []
      this.size = Math.ceil(this.total/this.itemsPerPage)

      if (this.size < this.step * 2 + 6) {
        // Case without any ellipse breaks
        this.add(1, this.size + 1);
      }
      else if (this.current < this.step * 2 + 1) {
        // Case with ellipse breaks at end
        this.add(1, this.step * 2 + 4);
        this.last();
      }
      else if (this.current > this.size - this.step * 2) {
        // Case with ellipse breaks at beginning
        this.first();
        this.add(this.size - this.step * 2 - 2, this.size + 1);
      }
      else {
        // Case with ellipse breaks at beginning and end
        this.first();
        this.add(this.current - this.step, this.current + this.step + 1);
        this.last();
      }
    }
  }
}
</script>
