/**
 * Tempoary no backend implementation.
 */

async function listPosts() {
  var res = await fetch(`/posts.json`);
  var data = await res.json();
  return data;
}

export default listPosts;
