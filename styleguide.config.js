const path = require('path');
module.exports = {
  require: [
    'babel-polyfill',
    path.join(__dirname, 'node_modules/uswds/dist/css/uswds.css'),
  ],
  components() {
    return ['src/shared/Alert/index.jsx'];
  },
};
