<template>
  <svg :width="getWidth" :height="getHeight" x="0.0" y="0.0" viewBox="0 0 48 48">
    <path
      fill="currentColor"
      fill-rule="nonzero"
      d="M24 .21L.21 24 24 47.79 47.79 24 24 .21zM5.82 24L24 5.82l7.69 7.69-2.8 2.8L24 11.43 11.43 24l4.88 4.89-2.8 2.8L5.82 24zm25.15 0L24 30.97 17.03 24 24 17.03 30.97 24zM16.31 34.49l2.8-2.8L24 36.57 36.57 24l-4.89-4.89 2.8-2.8 7.7 7.69L24 42.18l-7.69-7.69z"
    />
  </svg>
</template>

<script>
export default {
  props: {
    width: {
      type: Number,
      default: 40
    }
  },

  computed: {
    getWidth() {
      return this.width;
    },

    getHeight() {
      return Math.ceil((this.getWidth / 48) * 48);
    }
  }
};
</script>

