import 'babel-polyfill'
import 'whatwg-fetch'
import { doGenerator, template } from './util'
import { update, mkRoute } from './update'

const ROVERS = ['spirit', 'curiosity', 'opportunity']

// an object for app state
let State = {}

document.addEventListener('DOMContentLoaded', function() {
  // cache some dom elements
  State.main = document.querySelector('div.wrapper-main')
  State.nodes = Array.from(document.querySelectorAll('div.wrapper-item'))
  // load in the images
  let btnForward = document.getElementById('ctrl-forward')
  let btnBackward = document.getElementById('ctrl-backward')

  // show the first content item
  State.nodes[State.visible].classList.remove('hidden')

  // event handler where the magic happens
  let scrollHandler = (e) => {
    let up
    // for wheel events
    if (e.deltaY && e.deltaY < 0) up = true
    if (e.deltaY && e.deltaY > 0) up = false
    // for touch events
    if (e.type === 'touchmove') {
      if (State.lastX && State.lastX > e.touches[0].clientX) up = true
      if (State.lastX && State.lastX < e.touches[0].clientX) up = false
      State.lastX = e.touches[0].clientY
    }
    changeScene(up)
  }

  let btnDownHandler = (e) => {
    if (!e.target || !e.target.matches('button.ctrl-btn')) {
      return
    }
    let up = e.target.id === 'ctrl-backward' ? true : false
    changeScene(up)
  }

  ;(function init(){
    State.page = 1
    State.rover = 0
    State.tick = false
    State.visible = 0
    State.lastX = undefined
    State.fetched = []

    let route = mkRoute(ROVERS[State.rover], State.page)
    
    doGenerator(update, route)
    document.addEventListener('mousewheel', scrollHandler)
    document.addEventListener('DOMMouseScroll', scrollHandler)
    document.addEventListener('touchmove', scrollHandler)

    document.addEventListener('mousedown', btnDownHandler)
    document.addEventListener('mousedown', btnDownHandler)
  })()
})

/**
* Make new dom nodes with image data
* @param {Boolean} back - true if the previous image should be shown instead of the next image
*
* @return {Object}
*/
function changeScene(back) {
  if (back && State.visible === 0) return // don't do anything if someone scrolls up/clicks back right away
  if (!State.tick) {
    window.requestAnimationFrame(function() {
      // manage visible image state
      // first hide the current image
      State.nodes[State.visible].classList.add('hidden')
      // then figure out which one should be shown next
      if (!back) State.visible++
      if (back) State.visible--
      if (State.visible >= State.nodes.length) State.visible = State.nodes.length - 1
      if (State.visible < 0) State.visible = 0
      // now show the correct visible image
      State.nodes[State.visible].classList.remove('hidden')
      let route = makeRoute(ROVERS[State.rover], State.page)
      doGenerator(update, route)
      State.tick = false
    })

  }
  State.tick = true
}
