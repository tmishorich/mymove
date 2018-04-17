import React from 'react';
import PrivateRoute from 'shared/User/PrivateRoute';
import WizardPage from 'shared/WizardPage';

import Agreement from 'scenes/Legalese';
import Transition from 'scenes/Moves/Transition';
import MoveType from 'scenes/Moves/MoveTypeWizard';
import PpmSize from 'scenes/Moves/Ppm/PPMSizeWizard';
import PpmWeight from 'scenes/Moves/Ppm/Weight';
import SMNameWizard from 'scenes/ServiceMembers/SMNameWizard';
import ContactInfo from 'scenes/ServiceMembers/ContactInfo';

const Placeholder = props => {
  return (
    <WizardPage
      handleSubmit={() => undefined}
      pageList={props.pageList}
      pageKey={props.pageKey}
    >
      <div className="Todo">
        <h1>Placeholder for {props.title}</h1>
        <h2>{props.description}</h2>
      </div>
    </WizardPage>
  );
};

const stub = (key, pages, description) => ({ match }) => (
  <Placeholder
    pageList={pages}
    pageKey={key}
    title={key}
    description={description}
  />
);

const goHome = props => () => props.push('/');
const createMove = props => () => props.hasMove || props.createMove({});
const always = () => true;
const incompleteServiceMember = props => !props.hasCompleteProfile;
const hasHHG = ({ selectedMoveType }) =>
  selectedMoveType !== null && selectedMoveType !== 'PPM';
const hasPPM = ({ selectedMoveType }) =>
  selectedMoveType !== null && selectedMoveType !== 'HHG';
const isCombo = ({ selectedMoveType }) =>
  selectedMoveType !== null && selectedMoveType === 'COMBO';
const pages = {
  '/service-member/:serviceMemberId/create': {
    isInFlow: incompleteServiceMember,
    render: stub,
    description: 'Create your profile',
  },
  '/service-member/:serviceMemberId/name': {
    isInFlow: incompleteServiceMember,
    render: (key, pages) => ({ match }) => (
      <SMNameWizard pages={pages} pageKey={key} match={match} />
    ),
  },
  '/service-member/:serviceMemberId/contact-info': {
    isInFlow: incompleteServiceMember,
    render: (key, pages) => ({ match }) => (
      <ContactInfo pages={pages} pageKey={key} match={match} />
    ),
  },
  '/service-member/:serviceMemberId/duty-station': {
    isInFlow: incompleteServiceMember,
    render: stub,
    description: 'current duty station',
  },
  '/service-member/:serviceMemberId/residence-address': {
    isInFlow: incompleteServiceMember,
    render: stub,
    description: 'Current residence address',
  },
  '/service-member/:serviceMemberId/backup-mailing-address': {
    isInFlow: incompleteServiceMember,
    render: stub,
    description: 'Backup mailing address',
  },
  '/service-member/:serviceMemberId/backup-contacts': {
    isInFlow: incompleteServiceMember,
    render: stub,
    description: 'Backup contacts',
  },
  '/service-member/:serviceMemberId/transition': {
    isInFlow: incompleteServiceMember,
    render: stub,
    description: "OK, your profile's complete",
  },
  '/orders/:serviceMemberId/': {
    isInFlow: incompleteServiceMember,
    render: stub,
    description: 'Tell us about your move orders',
  },
  '/orders/:serviceMemberId/upload': {
    isInFlow: incompleteServiceMember,
    render: stub,
    description: 'Upload your orders',
  },
  '/orders/:serviceMemberId/complete': {
    isInFlow: incompleteServiceMember, //todo: this is probably not the right check
    render: (key, pages, description, props) => ({ match }) => {
      return (
        <WizardPage
          handleSubmit={createMove(props)}
          isAsync={!props.hasMove}
          hasSucceeded={props.hasMove}
          additionalParams={{ moveId: props.moveId }}
          pageList={pages}
          pageKey={key}
        >
          <div className="Todo">
            <h1>Placeholder for {key}</h1>
            <h2>Creating move here</h2>
          </div>
        </WizardPage>
      );
    },
  },
  '/moves/:moveId': {
    isInFlow: always,
    render: (key, pages) => ({ match }) => (
      <MoveType pages={pages} pageKey={key} match={match} />
    ),
  },
  '/moves/:moveId/schedule': {
    isInFlow: hasHHG,
    render: stub,
    description: 'Pick a move date',
  },
  '/moves/:moveId/address': {
    isInFlow: hasHHG,
    render: stub,
    description: 'enter your addresses',
  },

  '/moves/:moveId/ppm-transition': {
    isInFlow: isCombo,
    render: (key, pages) => ({ match }) => (
      <WizardPage handleSubmit={() => undefined} pageList={pages} pageKey={key}>
        <Transition />
      </WizardPage>
    ),
  },
  '/moves/:moveId/ppm-start': {
    isInFlow: state => state.selectedMoveType === 'PPM',
    render: stub,
    description: 'pickup zip, destination zip, secondary pickup, temp storage',
  },
  '/moves/:moveId/ppm-size': {
    isInFlow: hasPPM,
    render: (key, pages) => ({ match }) => (
      <PpmSize pages={pages} pageKey={key} match={match} />
    ),
  },
  '/moves/:moveId/ppm-incentive': {
    isInFlow: hasPPM,
    render: (key, pages) => ({ match }) => (
      <PpmWeight pages={pages} pageKey={key} match={match} />
    ),
  },
  '/moves/:moveId/review': {
    isInFlow: always,
    render: stub,
    description: 'Review',
  },
  '/moves/:moveId/agreement': {
    isInFlow: always,
    render: (key, pages, description, props) => ({ match }) => {
      return (
        <WizardPage handleSubmit={goHome(props)} pageList={pages} pageKey={key}>
          <Agreement match={match} />
        </WizardPage>
      );
    },
  },
};
export const getPageList = state =>
  Object.keys(pages).filter(pageKey => {
    const page = pages[pageKey];
    return page.isInFlow(state);
  });

export const getWorkflowRoutes = props => {
  const pageList = getPageList(props);
  return pageList.map(key => {
    const currPage = pages[key];
    const render = currPage.render(key, pageList, currPage.description, props);
    return <PrivateRoute exact path={key} key={key} render={render} />;
  });
};
