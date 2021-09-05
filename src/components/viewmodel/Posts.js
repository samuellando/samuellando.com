import Heading from '@/components/view/Heading.vue'
import ListPosts from '@/components/model/Posts.js'

export default {
  name: 'Posts',
  components: {
    Heading
  },
  data() {
    return {
      posts: [
      ]
    }
  },
  async mounted() {
    const posts = await ListPosts();
    this.posts = posts;
  }
}
