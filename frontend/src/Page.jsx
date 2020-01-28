import React from 'react';
import { Remarkable } from 'remarkable';

import PageNavigation from "./PageNavigation";

export default class Page extends React.Component {
  constructor(props) {
    super();
    this.state = {
      title: null,
      text: null,
      i: props.i
    };

    var callback = (state) => this.setState(state);
    this.i.getPage(props.id).then(
      function(response) {
        callback(response.data);
      },
      function(error) {
        callback({title: "Page Not Found."});
      }
    );
  }

  getMarkup() {
    var md = new Remarkable();
    return {__html: md.render(this.state.text)};
  }

  render() {
    return(
      <div id='pages'>
        <PageNavigation />
        <h1 onClick={() => edit(this.props.id)}>edit</h1>
        <h1>{this.state.title}</h1>
        <div id='page' dangerouslySetInnerHTML={this.getMarkup()}>
        </div>
      </div>
    );
  }
}