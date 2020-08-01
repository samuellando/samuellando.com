<template>
  <div class="sky">
    <!-- TODO separation of concerns. -->
    <Star v-for="(star, index) in createdStars" :key="index" :obj="star" />
  </div>
</template>

<script src="@/components/viewmodel/Sky.js"></script>
<style scoped src="@/components/view/stylesheets/Sky.css"></style>