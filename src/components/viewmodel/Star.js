export default {
  name: 'Star',

  data() {
    return {
      MAX_SIZE: 3,

      top: this.obj.top,
      left: this.obj.left,
      scale: this.obj.scale
    };
  },

  props: {
      obj: Object
  },

  computed: {
    style() {
      return {
        top: this.top+"vh",
        left: this.left+"vw",
        width: this.scale*this.MAX_SIZE+"px",
        height: this.scale*this.MAX_SIZE+"px"
      }
    }
  }
}