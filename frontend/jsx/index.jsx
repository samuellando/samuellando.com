class Signature extends React.Component {
  render() {
    return(
      <p onClick={redirectToGithub} className='signature'>
        Made with {this.props.emoji} by {this.props.by}
      </p>
    );
  }
}

class Navigation extends React.Component {
  render() {
    return(
      // TODO add planet control.
      <ul className='navigation'>
        <h2 onClick={loadPageMenu}>Pages</h2>
      </ul>
    );
  }
}

class Home extends React.Component {
  render() {
    return (
      <div id='index'>
        <Sky stars={1000} tick={20} />
        <h1 onClick={redirectToGithub}>Samuel Lando</h1>
        <System planets={1} />
        <Navigation />
        <Signature emoji="❤️" by="Sam" />
      </div>
    );
  }
}

function loadHome() {
  ReactDOM.render(
    <Home />,
    document.getElementById('root')
  );
}

function redirectToGithub() {
  window.location.href = "https://github.com/samuellando";
}

loadHome();
