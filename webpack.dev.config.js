var autoprefixer = require('autoprefixer')

module.exports = {
  entry: './frontend/index.js',
  output: {
    path: './build/assets',
    filename: 'app.js'
  },
  module: {
    loaders: [
      { test: /\.js$/, exclude: /node_modules/, loader: 'babel' },
      { test: /\.css$/, loader: 'style!css?modules&importLoaders=1&localIdentName=[name]__[local]___[hash:base64:5]!postcss' },
      { test: /\.(png|gif)$/, loader: 'url-loader?limit=8192' }
    ]
  },
  postcss: function() {
    return [autoprefixer]
  }
}
