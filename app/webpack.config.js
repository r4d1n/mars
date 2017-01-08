module.exports = {
  entry: ['babel-polyfill', 'whatwg-fetch', './js/index.js'],
  output: {
    path: './static',
    filename: 'bundle.js'
  },
  devtool: "inline-source-map",
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
