import React, { Component } from 'react';
import { Redirect, Route, Switch } from 'react-router-dom';
import { ConnectedRouter } from 'react-router-redux';
import { history } from 'shared/store';

import QueueList from './QueueList';

class QueueTable extends Component {
  render() {
    return (
      <div style={{ background: 'rgb(200,255,200)' }}>
        <h3>QueueTable</h3>
        <p>
          Now showing the <strong>{this.props.queueType}</strong> queue!
        </p>
      </div>
    );
  }
}

class QueueHeader extends Component {
  render() {
    return (
      <div style={{ background: 'rgb(200,200,255)' }}>
        <h1 style={{ margin: 0 }}>QueueHeader</h1>
      </div>
    );
  }
}

class Queues extends Component {
  render() {
    return (
      <div className="usa-grid">
        <div className="usa-width-one-fourth">
          <QueueList />
        </div>
        <div className="usa-width-three-fourths">
          <QueueTable queueType={this.props.match.params.queueType} />
        </div>
      </div>
    );
  }
}

class OfficeWrapper extends Component {
  componentDidMount() {
    document.title = 'Transcom PPP: Office';
  }

  render() {
    return (
      <ConnectedRouter history={history}>
        <div className="Office site">
          <main className="site__content">
            <div>
              <div className="usa-grid">
                <QueueHeader />
              </div>
              <Switch>
                <Redirect from="/" to="/queues/new_moves" exact />
                <Route path="/queues/:queueType" component={Queues} />
              </Switch>
            </div>
          </main>
        </div>
      </ConnectedRouter>
    );
  }
}

export default OfficeWrapper;
