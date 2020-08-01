import StarModel from '@/components/model/Star.js'

import Star from "@/components/view/Star.vue"

export default {
  name: 'Sky',

  data() {
    return {
      TICK_RATE: 1,
      ticker: null,
      createdStars: []
    };
  },

  props: {
    stars: Number
  },

  computed: {
    
  },

  mounted() {
    this.startTick();
  },

  methods: {
    startTick() {
      this.ticker = setInterval(this.tick, this.TICK_RATE);
    },
    stopTick() {
      clearInterval(this.ticker);
    },
    tick() {
      if (this.createdStars.length >= this.stars) {
        this.stopTick();
      } else {
        this.addNewStar();
      }
    },
    addNewStar() {
      var star = StarModel.factory(
        Math.random() * 100.0,
        Math.random() * 100.0,
        Math.random()
      );
      this.createdStars.push(star)
    },
  },

  components: {
    Star
  }
}