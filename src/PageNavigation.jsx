import React from 'react';

import loadHome from './App';
import loadPageMenu from './App';

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