class Star extends React.Component {
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

class Sky extends React.Component {
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
    return (
      <div className='sky'>
        {this.state.stars.map(star => <Star top={star.top} right={star.right} size={star.size} />)}
      </div>
    );
  }
}
