<template>
  <ul class="links">
    <li><a href="https://github.com/samuellando"><img src="/github.png" width="75"/></a></li>
    <li><a href="https://www.linkedin.com/in/samuel-lando-7a71a2155/"><img src="/linkedin.png" width="75"/></a></li>
    <li><a id="internal-link" :href="link">{{text}}</a></li>
    <li :style="{'visibility':info}"><a href="#info"><img src="/arrow.png" width="75" class="animate-flicker"></a></li>
  </ul>
</template>

<script src="@/components/viewmodel/Links.js"></script>
<style scoped src="@/components/view/stylesheets/Links.css"></style>
