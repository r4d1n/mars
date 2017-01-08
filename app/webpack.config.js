module.exports = {
  entry: {
    init: './js/init.js',
    bundle: ['babel-polyfill', 'whatwg-fetch', './js/index.js']
  },
  output: {
    path: './static/dist',
    filename: '[name].js'
  },
  devtool: 'inline-source-map',
  module: {
    loaders: [
      {
        test: /\.js$/,
        loader: 'babel-loader',
        query: {
          presets: ['es2015']
        }
      }
    ]
  },
  resolve: {
    extensions: ['', '.js', '.json']
  }
}
