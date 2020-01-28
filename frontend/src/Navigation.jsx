import React from 'react';

export default class Navigation extends React.Component {
  render() {
    return(
      // TODO add planet control.
      <ul className='navigation'>
        <h2 onClick={loadPageMenu}>Pages</h2>
      </ul>
    );
  }
}