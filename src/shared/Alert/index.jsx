// eslint-disable-next-line no-unused-vars
import React from 'react';
import PropTypes from 'prop-types';

//this is taken from https://designsystem.digital.gov/components/alerts/
const Alert = props => (
  <div className={`usa-alert usa-alert-${props.type}`}>
    <div className="usa-alert-body">
      <h3 className="usa-alert-heading">{props.heading}</h3>
      <p className="usa-alert-text">{props.children}</p>
    </div>
  </div>
);

Alert.propTypes = {
  heading: PropTypes.string.isRequired,
  children: PropTypes.node.isRequired,
  type: PropTypes.oneOf(['error', 'warning', 'info', 'success']),
};
export default Alert;
