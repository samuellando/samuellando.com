import Sky from '@/components/view/Sky.vue'
import Heading from '@/components/view/Heading.vue'
import Links from '@/components/view/Links.vue'
import Signature from '@/components/view/Signature.vue'
import Posts from '@/components/view/Posts.vue'
import Post from '@/components/view/Post.vue'

export default {
  name: 'Blog',
  components: {
    Sky,
    Heading,
    Links,
    Posts,
    Post,
    Signature,
  },

  data() {
    var path = window.location.pathname;
    console.log(path);
    if (path == "/blog")
      return {
        type: "feed",
        pillhome: false,
        pillblog: true,
        noback: true,
        backlink: ""
      };
    else
      return {
        type: "post",
        id: 1,
        pillhome: false,
        pillblog: false,
        noback: false,
        backlink: "/blog"
      };
  }
}
