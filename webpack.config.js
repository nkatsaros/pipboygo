var autoprefixer = require('autoprefixer')
var webpack = require('webpack')

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
  },
  plugins: [
    new webpack.DefinePlugin({
      'process.env': {
          NODE_ENV: JSON.stringify('production')
      }
    }),
    new webpack.optimize.UglifyJsPlugin({
        compress: {
            warnings: false
        }
    })
  ]
}
