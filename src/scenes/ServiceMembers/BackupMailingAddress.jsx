import { pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { updateServiceMember, loadServiceMember } from './ducks';
import { reduxifyForm } from 'shared/JsonSchemaForm';
import { no_op } from 'shared/utils';
import WizardPage from 'shared/WizardPage';

const subsetOfFields = ['backup_mailing_address'];

const uiSchema = {
  title: 'Backup mailing address',
  description:
    'Enter a backup mailing address, such as your permanent residence or a parent’s address.',
  order: subsetOfFields,
  definitions: {
    Address: {
      order: [
        'street_address_1',
        'street_address_2',
        'city',
        'state',
        'postal_code',
      ],
    },
  },
  requiredFields: subsetOfFields,
};
const formName = 'service_member_backup_mailing_addresss';
const CurrentForm = reduxifyForm(formName);

export class BackupMailingAddress extends Component {
  componentDidMount() {
    this.props.loadServiceMember(this.props.match.params.serviceMemberId);
  }

  handleSubmit = () => {
    const pendingValues = this.props.formData.values;
    if (pendingValues) {
      const patch = pick(pendingValues, subsetOfFields);
      this.props.updateServiceMember(patch);
    }
  };

  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      currentServiceMember,
    } = this.props;
    const isValid = this.refs.currentForm && this.refs.currentForm.valid;
    const isDirty = this.refs.currentForm && this.refs.currentForm.dirty;
    // initialValues has to be null until there are values from the action since only the first values are taken
    const initialValues = currentServiceMember
      ? pick(currentServiceMember, subsetOfFields)
      : null;
    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        isAsync={true}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={isValid}
        pageIsDirty={isDirty}
        hasSucceeded={hasSubmitSuccess}
        error={error}
      >
        <CurrentForm
          ref="currentForm"
          className={formName}
          handleSubmit={no_op}
          schema={this.props.schema}
          uiSchema={uiSchema}
          initialValues={initialValues}
        />
      </WizardPage>
    );
  }
}
BackupMailingAddress.propTypes = {
  schema: PropTypes.object.isRequired,
  updateServiceMember: PropTypes.func.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { updateServiceMember, loadServiceMember },
    dispatch,
  );
}
function mapStateToProps(state) {
  const props = {
    schema: {},
    formData: state.form[formName],
    ...state.serviceMember,
  };
  if (state.swagger.spec) {
    props.schema = state.swagger.spec.definitions.CreateServiceMemberPayload;
  }
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(
  BackupMailingAddress,
);
