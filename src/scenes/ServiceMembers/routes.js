// TODO: this file is obsolete--all changes should go in shared/AppWrapper/getWorkflowRoutes

import React from 'react';

import PrivateRoute from 'shared/User/PrivateRoute';
import WizardPage from 'shared/WizardPage';
import SMNameWizard from 'scenes/ServiceMembers/SMNameWizard';

const Placeholder = props => {
  return (
    <WizardPage
      handleSubmit={() => undefined}
      pageList={props.pageList}
      pageKey={props.pageKey}
    >
      <h1>Placeholder for {props.title}</h1>
    </WizardPage>
  );
};

const stub = (key, pages, component) => ({ match }) => {
  if (component) {
    const pageComponent = React.createElement(component, { match }, null);
    return (
      <WizardPage handleSubmit={() => undefined} pageList={pages} pageKey={key}>
        {pageComponent}
      </WizardPage>
    );
  } else {
    return <Placeholder pageList={pages} pageKey={key} title={key} />;
  }
};

export default () => {
  const pages = {
    '/service-member/:serviceMemberId/create': { render: stub },
    '/service-member/:serviceMemberId/name': {
      render: (key, pages) => ({ match }) => (
        <SMNameWizard pages={pages} pageKey={key} match={match} />
      ),
    },
    '/service-member/:serviceMemberId/contact-info': { render: stub },
    '/service-member/:serviceMemberId/duty-station': { render: stub },
    '/service-member/:serviceMemberId/residence-address': {
      render: stub,
    },
    '/service-member/:serviceMemberId/backup-mailing-address': { render: stub },
    '/service-member/:serviceMemberId/backup-contacts': { render: stub },
    '/service-member/:serviceMemberId/transition': { render: stub },
  };
  const pageList = Object.keys(pages);
  const componentMap = {};
  return pageList.map(key => {
    const step = key.split('/').pop();
    var component = componentMap[step];
    const render = pages[key].render(key, pageList, component);
    return <PrivateRoute exact path={key} key={key} render={render} />;
  });
};