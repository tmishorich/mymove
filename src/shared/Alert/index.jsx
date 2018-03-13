// eslint-disable-next-line no-unused-vars
import React from 'react';
import PropTypes from 'prop-types';

/** An Alert component that produces  https://designsystem.digital.gov/components/alerts/
 */
const Alert = props => (
  <div className={`usa-alert usa-alert-${props.type}`}>
    <div className="usa-alert-body">
      <h3 className="usa-alert-heading">{props.heading}</h3>
      <p className="usa-alert-text">{props.children}</p>
    </div>
  </div>
);

Alert.propTypes = {
  /** the heading of the Alert */
  heading: PropTypes.string.isRequired,
  /** the type of the Alert */
  type: PropTypes.oneOf(['error', 'warning', 'info', 'success']).isRequired,
  /** the content of the Alert */
  children: PropTypes.node.isRequired,
};
export default Alert;
