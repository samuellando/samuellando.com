<template>
  <div class="card m-5">
    <div class="card-body">
      <h1 class="card-title">{{title}}</h1>
      <img :src="image" />
      <h3 class="card-title">{{description}}</h3>
      <div class="card-text text-start m-5"><span v-html="content"></span></div>
    </div>
  </div>
</template>

<script src="@/components/viewmodel/Post.js"></script>
<style scoped src="@/components/view/stylesheets/Post.css"></style>
