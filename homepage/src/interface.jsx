

var backendURL = "http://localhost:8080";

export default class Interface {
  getPages() {
    return ajax("GET", backendURL+"/pages")
  }

  getPage(id) {
    return ajax("GET", backendURL+"/page/"+id)
  }

  createPage(title, text) {
    return ajax("POST", backendURL+"/page", {title:title, text:text});
  }

  updatePage(id, title, text) {
    return ajax("PUT", backendURL+"/page/"+id, {title:title, text:text});
  }

  deletePage(id) {
    return ajax("DELETE", backendURL+"/page/"+id);
  }
}

function ajax(method, url, payload) {
  return new Promise(function(resolve, reject) {
    var req = new XMLHttpRequest();
    req.open(method, url);
    req.onreadystatechange = function() {
      if (req.readyState == 4) {
        if (req.status == 200) {
          var res = {
            status: req.status,
            data: JSON.parse(req.responseText)
          };
          resolve(res);
        } else {
          reject({status: req.status});
        } 
      }
    };
    req.send(JSON.stringify(payload));
  });
}
