export default {
  name: 'links',
  props: {
    links: Array
  },
  data() {
    if (window.location.pathname == "/")
      return {
        link: "/blog",
        text: "BLOG",
        info: "visible"
      }
    else
      return {
        link: "/",
        text: "HOME",
        info: "hidden"
      }
  }
};
