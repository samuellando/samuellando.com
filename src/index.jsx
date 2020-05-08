import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';

export default function loadApp() {
  ReactDOM.render(
    <App />,
    document.getElementById('root')
  );
}

loadApp();