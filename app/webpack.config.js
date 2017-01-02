module.exports = {
  entry: ['babel-polyfill', './javascript/index.js'],
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
          presets: ['es2015', 'react']
        }
      },
      {
        test: /\.scss$/,
        loaders: ["style-loader", "css-loader?sourceMap", "sass-loader?sourceMap"]
      }
    ]
  },
  resolve: {
    extensions: ['', '.js', '.json']
  }
};
