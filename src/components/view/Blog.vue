<template>
  <div id="app">
    <b-nav pills class="m-2">
      <b-nav-item :active="pillhome" href="/">Home</b-nav-item>
      <b-nav-item :active="pillblog" href="/blog">Blog</b-nav-item>
      <b-nav-item :disabled="noback" :href="backlink">Back</b-nav-item>
    </b-nav>
    <div class="blogContainer">
      <Sky :stars="500" />
      <Heading text="Samuel Lando / Blog" />
      <div v-if="type == 'feed'">
        <Posts />
      </div>
      <div v-else>
        <Post />
      </div>
      <Signature text="Made with ♥ by Sam"/>
    </div>
  </div>
</template>

<script src="@/components/viewmodel/Blog.js"></script>
<style src="@/components/view/stylesheets/Blog.css"></style>
