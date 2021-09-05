/**
 * Tempoary no backend implementation.
 */

import Marked from "marked"

async function getPost(id) {
  var res = await fetch(`/${id}.md`);
  var md = await res.text();
  res  = await fetch(`/${id}.json`);
  var data = await res.json();
  return {
    title: data.title,
    description: data.description,
    content: Marked(md),
    image: data.image,
    listed: data.listed
  };
}

export default getPost;
