class Signature extends React.Component {
  render() {
    return(
      <p className='signature'>
        Made with {this.props.emoji} by {this.props.by}
      </p>
    );
  }
}

class Navigation extends React.Component {
  render() {
    return(
      <ul className='navigation'>
        <li><h2>Planet Control</h2></li>
        <li><h2>Pages</h2></li>
      </ul>
    );
  }
}

class Home extends React.Component {
  render() {
    return (
      <div>
        <Sky stars={1000} tick={20} />
        <h1>Samuel Lando</h1>
        <System planets={1} />
        <Navigation />
        <Signature emoji="❤️" by="Sam" />
      </div>
    );
  }
}

ReactDOM.render(
  <Home />,
  document.getElementById('root')
);
