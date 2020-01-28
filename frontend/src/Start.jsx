import React from 'react';

export default class Star extends React.Component {
  render() {
        var style = {
      position: 'absolute',
      top: this.props.top+"vh",
      left: this.props.right+"vw",
      height: this.props.size+"px",
      width: this.props.size+"px",
    }

    return(
      <div className='star' style={style}>
      </div>
    );
  }
}