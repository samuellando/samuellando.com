import React from 'react';
import ReactDOM from 'react-dom';
import { Remarkable } from 'remarkable';

import loadHome from './index';

import Interface from './interface';

class Page extends React.Component {
  constructor(props) {
    super();
    this.state = {
      title: null,
      text: null
    };
    this.i = new Interface();
    var callback = (state) => this.setState(state);
    this.i.getPage(props.id).then(
      function(response) {
        callback(response.data);
      },
      function(error) {
        callback({title: "Page Not Found."});
      }
    );
  }

  getMarkup() {
    var md = new Remarkable();
    return {__html: md.render(this.state.text)};
  }

  render() {
    return(
      <div id='pages'>
        <PageNavigation />
        <h1>{this.state.title}</h1>
        <div id='page' dangerouslySetInnerHTML={this.getMarkup()}>
        </div>
      </div>
    );
  }
}

class PageNavigation extends React.Component {
  render() {
    return(
      <ul>
        <li onClick={loadHome}><h2>Home</h2></li>
        <li onClick={loadPageMenu}><h2>Menu</h2></li>
      </ul>
    );
  }
}

class PageMenu extends React.Component {
  constructor() {
    super();
    this.state = {pages: []};
    this.i = new Interface();
    var callback = (state) => this.setState(state);
    this.i.getPages().then(
      function(response) {
        callback({pages: response.data});
      },
      function(error) {
        callback({pages: null});
      }
    );
  }

  render() {
    var output;
    if (this.state.pages == null) {
      output = <h1>No Pages Found.</h1>;
    } else if (this.state.pages.length == 0) {
      output = <h1>Loading...</h1>;
    } else {
      output = <ul>
          {this.state.pages.map(
            page => <li onClick={() => loadPage(page.id)}><h2>{page.title}</h2></li>
          )}
        </ul>;
    }
    return(
      <div id='pages'>
        <PageNavigation />
        {output}
      </div>
    );
  }
}

export default function loadPageMenu() {
  ReactDOM.render(<PageMenu />, document.getElementById('root'));
}

function loadPage(id) {
  ReactDOM.render(<Page id={id} />, document.getElementById('root'));
}