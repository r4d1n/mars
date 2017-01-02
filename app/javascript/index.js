"use strict";

require("../styles/index.scss")

import "babel-polyfill"

document.addEventListener("DOMContentLoaded", function() {
  // cache some dom elements
  let wrappers = document.querySelectorAll("div.wrapper-item")
  let photos = document.querySelectorAll("img.photo")
  // set up some globals
  const routeFn = (rover, page) => `/rover/${rover}/page/${page}`
  let page = 1
  let rover = "curiosity"

  ;(function init() {
    photos = doLazyLoad(photos)
    let route = routeFn(rover, page)
    console.log("initialize", page, rover, route)
    // Util.runGenerator(updateImages, route)
  })()

  /**
  * Lazy load initial images in rendered markup
  * @param {nodelist} imgNodes - an array-like list of <img> DOM nodes in the initial index.html
  */
  function doLazyLoad(imgNodes) {
    console.log("in lazy load", imgNodes);
    return Array.prototype.map.call(imgNodes, (el) => el.src = el.dataset.src)
  }

  /**
  * Fetch images and metadata and append to the DOM
  * @param {string} uri - a generator function to run asynchronously
  */
  function* updateImages(uri) {
    try {
      let res = yield fetch(uri) // returns a promise for the response
      let list = yield res.json() // returns a promise for json
      let i = 0
      let j = 0
      // while (i < divs.)
      // console.log(el)
      // do dom updates
      // let div = document.createElement("div")
      // let img = document.createElement("img")
      // img.src = el.img_src
      // div.appendChild(img)
      // document.body.appendChild(div)
    } catch (err) {
      console.error(err)
    }
  }
})
