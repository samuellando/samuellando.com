<template>
  <div id="app">
    <Sky :stars="500" />
    <Heading text="Samuel Lando" />
    <System :planets="3" />
    <Links/>
    <Signature msg="Made with love by Sam"/>
  </div>
</template>

<script src="@/components/viewmodel/App.js"></script>
<style src="@/components/view/stylesheets/App.css"></style>