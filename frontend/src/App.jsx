import React from 'react';

import Navigation from "./Navigation";
import Signature from "./Signature";
import Sky from './Sky';
import System from './System';

import './css/App.css';

export default class App extends React.Component {
  constructor(props) {
    super()
    this.state = {
      mode: "home"
    }
  }

  loadHome() {
    this.setState({mode: "home"});
  }

  loadPageMenu() {
    this.setState({mode: "pageMenu"});
  }

  render() {
    var out = "";
    if (this.state.mode === "home") {
      out = <div id='index'>
        <Sky stars={1000} tick={20} />
        <h1 onClick={redirectToGithub}>Samuel Lando</h1>
        <System planets={1} />
        <Navigation app={this} />
        <Signature emoji="❤️" by="Sam" onClick={redirectToGithub} />
      </div>;
    } else if (this.state.mode === "pageMenu") {
      out = <h1>PageMenu</h1>;
    }
    return (out);
  }
}

function redirectToGithub() {
  window.location.href = "https://github.com/samuellando";
}