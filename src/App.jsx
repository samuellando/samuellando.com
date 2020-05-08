import React from 'react';

import Navigation from "./Navigation";
import Signature from "./Signature";
import Sky from './Sky';
import System from './System';
import Login from './Login';

import './css/App.css';

export default class App extends React.Component {
  constructor(props) {
    super()
    this.state = {
      mode: "home"
    }

    this.loadLogIn = this.loadLogIn.bind(this); 
    this.loadPageMenu = this.loadPageMenu.bind(this);
  }

  loadHome() {
    this.setState({mode: "home"});
  }

  loadPageMenu() {
    console.log(this);
    this.setState({mode: "pageMenu"});
  }
  
  loadLogIn() {
    console.log("OK");
    this.setState({mode: "logIn"});
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
    } else if (this.state.mode === "logIn") {
      out = <Login />;
    }
    return (out);
  }
}

function redirectToGithub() {
  window.location.href = "https://github.com/samuellando";
}