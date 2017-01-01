"use strict";


(function init(){
  let page = 0
  let rover = "curiosity"
  let route = `/rover/${rover}/page/${page}`
  console.log(page, rover, route)
  let divs = document.querySelectorAll("div.wrapper-item")
  UTIL.runGenerator(getImages, route)
}())


/**
* Fetch images and metadata and append to the DOM
* @param {string} uri - a generator function to run asynchronously
*/
function* getImages(uri) {
  try {
    let res = yield fetch(uri) // returns a promise for the response
    let list = yield res.json() // returns a promise for json
    list.forEach((el) => {
      console.log(el)
      // do dom updates
      // let div = document.createElement("div")
      // let img = document.createElement("img")
      // img.src = el.img_src
      // div.appendChild(img)
      // document.body.appendChild(div)
    })
  } catch (err) {
    console.error(err)
  }
}
