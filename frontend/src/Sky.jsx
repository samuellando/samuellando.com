import React from 'react';
import Star from './Star';

import './css/sky.css';

export default class Sky extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      stars: [],
      n: props.stars,
      tick: props.tick
    };
  }

  tick() {
    var maxSize = 2;
    var star = {top: 0, right: 0, size: 0};

    star.top = Math.random() * 100.0;
    star.right = Math.random() * 100.0;
    star.size = maxSize*Math.random();

    this.setState(() => ({stars: this.state.stars.concat([star])}));
    if (this.state.stars.length >= this.state.n) {
      clearInterval(this.interval);
    }
  }

  componentDidMount() {
    this.interval = setInterval(() => this.tick(), this.state.tick);
  }

  componentWillUnmount() {
    clearInterval(this.interval);
  }

  render() {
    var style = {
      zIndex: -1000000,
      position: "relative",
    };
    return (
      <div className='sky' style={style}>
        {this.state.stars.map(star => <Star top={star.top} right={star.right} size={star.size} />)}
      </div>
    );
  }
}
