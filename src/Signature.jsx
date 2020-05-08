import React from 'react';

export default class Signature extends React.Component {
  render() {
    return(
      <p onClick={this.props.onClick} className='signature'>
        Made with {this.props.emoji} by {this.props.by}
      </p>
    );
  }
}