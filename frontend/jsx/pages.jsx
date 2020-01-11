var pages = [
      {
        id: 0,
        title: 'downloads',
        text: '# dowloads page'
      },
      {
        id: 1,
        title: 'test1',
        text: '# test1 page'
      },
      {
        id: 2,
        title: 'test2',
        text: '# test2 page'
      },
    ];

class Page extends React.Component {
  constructor(props) {
    super();
    // TODO Fetch the mardown from the backend.
    var title = pages[props.id]['title'];
    var text = pages[props.id]['text'];

    this.state = {
      title: title,
      text: text,
    };
  }

  render() {
      return(
        this.state.text
      );
  }
}

class PageMenu extends React.Component {
  constructor() {
    super();
    // TODO fetch list of pages {id, title} from the backend.
    this.state = {
      pages: pages
    };
  }

  render() {
    console.log(this.state.pages);
    return(
      <div>
        <div id='page'>
          <ul>
            {this.state.pages.map(
              page => <li onClick={() => loadPage(page.id)}><h2>{page.title}</h2></li>
            )}
          </ul>
        </div>
      </div>
    );
  }
}

function loadPageMenu() {
  ReactDOM.render(<PageMenu />, document.getElementById('root'));
}

function loadPage(id) {
  ReactDOM.render(<Page id={id} />, document.getElementById('page'));
}
