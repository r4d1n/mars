'use strict';

// an object for app state
let State = {}
State.page = 0
State.rover = 'curiosity'
State.tick = false
State.queue = []
State.nodes = []
State.visibleIndex = 0

document.addEventListener('DOMContentLoaded', function() {
  // cache some dom elements
  let main = document.querySelector('div.wrapper-main')
  let items = Array.from(document.querySelectorAll('div.wrapper-item'))
  let photos = Array.from(document.querySelectorAll('img.photo'), (img) => Util.lazyLoad(img))

  // show the first content item
  items[State.visibleIndex].classList.remove('hidden')
  State.nodes.concat(photos)
  console.log(State)

  // load in the images

  let scrollHandler = (e) => {
    if (!State.tick) {
      window.requestAnimationFrame(function() {
        // let route = Util.makeRoute(State.rover, State.page)
        // Util.generator(updateImages, route)
        console.log('bounce bounce')
        items[State.visibleIndex].classList.add('hidden')
        State.visibleIndex++
        if (State.visibleIndex >= items.length) State.visibleIndex = 0
        items[State.visibleIndex].classList.remove('hidden')
        State.tick = false
      })
    }
    State.tick = true
  }

  ;(function init(){
    main.addEventListener('mousewheel', scrollHandler)
    main.addEventListener('DOMMouseScroll', scrollHandler)
    main.addEventListener('touchmove', scrollHandler)
  })()
})

/**
* Fetch images and metadata and append to the DOM
* @param {string} uri - a generator function to run asynchronously
*/
function* updateImages(uri) {
  try {
    let res = yield fetch(uri) // returns a promise for the response
    let list = yield res.json() // returns a promise for json
    State.page++
    for (let item of list) {
      State.queue.unshift(item)
      console.log(State.queue)
    }
  } catch (err) {
    console.error(err)
  }
}

function newSrc(node, src) {

}
