import React from 'react';

export default class PageEditor extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      title: 'Loading...',
      text: 'Loading...',
    };
    this.handleTitleChange = this.handleTitleChange.bind(this);
    this.handleTextChange = this.handleTextChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleDelete = this.handleDelete.bind(this);

    this.i = new Interface();
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

  handleTitleChange(event) {
    this.setState({title: event.target.value});
  }

  handleTextChange(event) {
    this.setState({text: event.target.value});
  }

  handleSubmit(event) {
    this.i = new Interface();
    if (this.props.id >= 0) {
      this.i.updatePage(this.props.id, this.state).then(
        function(response) {
          alert("Saved");
        },
        function(error) {
          alert("Not Able to Save");
        }
      );
    } else {
      this.i.createPage(this.state).then(
        function(response) {
          alert("Created");
        },
        function(error) {
          console.log(error);
          alert("Not Able to Create");
        }
      );
    }
    event.preventDefault();
  }

  handleDelete(event) {
    this.i = new Interface();
    if (this.props.id >= 0) {
      this.i.deletePage(this.props.id);
    }
    event.preventDefault();
    loadPageMenu();
  }

  render() {
    return (
      <div>
        <PageNavigation />
        <form onSubmit={this.handleSubmit}>
          <input type='text' value={this.state.title} onChange={this.handleTitleChange}/>
          <br/>
          <textarea value={this.state.text} onChange={this.handleTextChange} rows="4" cols="50">
          </textarea>
          <br/>
          <input type='submit' value='Save'/>
          <input type='button' value='Delete' onClick={this.handleDelete} />
        </form>
      </div>
    );
  }
}