var autoprefixer = require('autoprefixer');

module.exports = {
  entry: './frontend/index.js',
  output: {
    path: './build/assets',
    filename: 'app.js'
  },
  module: {
    loaders: [
      { test: /\.js$/, exclude: /node_modules/, loader: 'babel' },
      { test: /\.css$/, loader: 'style!css!postcss' }
    ]
  },
  postcss: function() {
    return [autoprefixer];
  }
};
