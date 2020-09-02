<template>
<div class="system" >
  <Body
    :radius="200"
    :tilt="0"
    :spin="0.001"
    :numberOfArcs="10"
    :arcTiltIncrement="2"
    arcBorderWidth="thin"
    arcBorderColor="yellow"
    arcFillColor="rgba(235, 229, 52, 0.10)"
    :arcFillOpacity="25"

    :eccentricity="0"
    :semimajorAxis="0"
    :inclination="0"
    :longitudeOfAscendingNode="0"
    :argumentOfPeriapsis="0"
    :trueAnomaly="0"
    :mobile="true"
  />
  <Body
    :radius="75"
    :tilt="0"
    :spin="0.005"
    :numberOfArcs="10"
    :arcTiltIncrement="2"
    arcBorderWidth="thin"
    arcBorderColor="orange"
    arcFillColor="rgba(235, 229, 52, 0.10)"
    :arcFillOpacity="25"

    :eccentricity="0"
    :semimajorAxis="0"
    :inclination="0"
    :longitudeOfAscendingNode="0"
    :argumentOfPeriapsis="0"
    :trueAnomaly="0"
    :mobile="true"
  />
  <Body
    :radius="100"
    :tilt="0.75"
    :spin="0.005"
    :numberOfArcs="10"
    :arcTiltIncrement="0"
    arcBorderWidth="thin"
    arcBorderColor="rgba(52, 235, 168, 1)"
    arcFillColor="rgba(52, 235, 168, 0.05)"
    :arcFillOpacity="100"

    :eccentricity="0.2"
    :semimajorAxis="200"
    :inclination="0"
    :longitudeOfAscendingNode="0.75"
    :argumentOfPeriapsis="0.75"
    :trueAnomaly="0"
    :mobile="false"
  />
  <Body
    :radius="100"
    :tilt="0"
    :spin="0.01"
    :numberOfArcs="15"
    :arcTiltIncrement="0.5"
    arcBorderWidth="thin"
    arcBorderColor="rgba(217, 52, 235, 0.5)"
    arcFillColor="rgba(52, 198, 235, 0.05)"
    :arcFillOpacity="100"

    :eccentricity="0.5"
    :semimajorAxis="300"
    :inclination="0"
    :longitudeOfAscendingNode="2.25"
    :argumentOfPeriapsis="0.75"
    :trueAnomaly="0"
    :mobile="false"
  />
</div>
</template>

<script src="@/components/viewmodel/System.js"></script>
<style scoped src="@/components/view/stylesheets/System.css"></style>