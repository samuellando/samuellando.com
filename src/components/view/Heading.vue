<template>
  <div class="heading">
    <h1>{{ msg }}</h1>
  </div>
</template>

<script>
export default {
  name: 'heading',
  props: {
    msg: String
  }
}
</script>

<style scoped>
</style>
