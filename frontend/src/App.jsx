import React from 'react';

import Navigation from "./Navigation";
import Signature from "./Signature";
import Sky from './Sky';
import System from './System';

import './css/App.css';

export default class App extends React.Component {
  render() {
    return (
      <div id='index'>
        <Sky stars={1000} tick={20} />
        <h1 onClick={redirectToGithub}>Samuel Lando</h1>
        <System planets={1} />
        <Navigation />
        <Signature emoji="❤️" by="Sam" onClick={redirectToGithub} />
      </div>
    );
  }
}

function redirectToGithub() {
  window.location.href = "https://github.com/samuellando";
}