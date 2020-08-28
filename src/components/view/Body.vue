<template>
<div class="body" :style="[style]" >
    <div v-for="(arcStyle, index) in arcStyles" :key="index" :style="[arcStyle]" >
    </div>
</div>
</template>

<script src="@/components/viewmodel/Body.js"></script>
<style scoped src="@/components/view/stylesheets/Body.css"></style>