<template>
  <div class="row row-cols-auto justify-content-center">
    <a v-for="post in posts" :key="post.id" :href="'/blog/'+post.id">
      <b-card
        :title="post.title"
        :img-src="post.image"
        img-alt="Alt"
        img-top
        tag="article"
        style="max-width: 40rem;"
        class="m-2"
      >
        <b-card-text>
          {{post.description}}
        </b-card-text>
      </b-card>
    </a>
  </div>
</template>

<script src="@/components/viewmodel/Posts.js"></script>
<style scoped src="@/components/view/stylesheets/Posts.css"></style>
