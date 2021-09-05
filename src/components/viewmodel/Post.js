import Heading from '@/components/view/Heading.vue'
import GetPost from '@/components/model/Post.js'

export default {
  name: 'Post',
  components: {
    Heading
  },
  data() {
    return {
      id: null,
      listed: null,
      title: null,
      description: null,
      image: null,
      content: null
    };
  },
  async mounted () {
    const data = await GetPost(window.location.pathname.split("/")[2]);
    for (var key in data) 
      this[key] = data[key];
  }
}
