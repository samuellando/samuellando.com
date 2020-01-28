import React from 'react';

import PageNavigation from './PageNavigation';
import loadPage from './App';

class PageMenu extends React.Component {
  constructor() {
    super();
    this.state = {
        pages: [],
        i: this.props.i
    };

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