const webpack = require("webpack")

module.exports = {
  entry: {
    bundle: ["whatwg-fetch", "./js/index.js"]
  },
  output: {
    path: "./dist/static",
    filename: "[name].js"
  },
  devtool: "inline-source-map",
  module: {
    loaders: [
      {
        test: /\.js$/,
        loader: "babel-loader",
        query: {
          presets: ["es2015"]
        }
      }
    ]
  },
  resolve: {
    extensions: ["", ".js", ".json"]
  }
  //, plugins: [
  //   new webpack.optimize.UglifyJsPlugin({
  //     compress: {
  //       warnings: false
  //     }
  //   })
  // ]
}
