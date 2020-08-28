<template>
  <div id="app">
    <div class="homeContainer">
      <Sky :stars="500" />
      <Heading text="Samuel Lando" />
      <System :planets="3" />
      <Links/>
      <Signature text="Made with ♥ by Sam"/>
    </div>
    <Information/>
  </div>
</template>

<script src="@/components/viewmodel/App.js"></script>
<style src="@/components/view/stylesheets/App.css"></style>