class Page extends React.Componenet {
  constructor(props) {
    // TODO Fetch the mardown from the backend.
    this.state = {
      title: title,
      text: text,
    };
  }

  render() {
      return(
        <h1>{this.state.title}</h1>
        //{markdownRender(this.state.text)}
      );
  }
}

class PageMenu extends React.Component {
  constructor() {
    // TODO fetch list of pages {id, title} from the backend.
    this.state = {
      pages: pages
    }
  }

  render() {
    return(
      <div>
        <PageNavigation />
        <div id='page'>
          <ul>
            {this.state.pages.map(
              //page => <li onClick={loadPage(page.id)}>page.title</li>
            )}
          </ul>
        </div>
      </div>
    );
  }
}

function loadPage(id) {
  reactDOM(<Page id={id} />, document.getElementById('page'));
}
