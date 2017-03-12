
// const ROVERS = ["curiosity", "spirit", "opportunity"]

// let Config = {}
// Config.limit = 10

// // an object for app state
// let State = {}

import { InfiniteScroll } from './infinitescroll';

document.addEventListener('DOMContentLoaded', function init() {
  // start lazy loading the first images right away
  let photos = Array.from(document.querySelectorAll('img.photo'), (img) => lazyLoad(img));

  let target = document.getElementById('scroll-target');
  console.log('target:', target);
  let scroll = new InfiniteScroll(target);

  // document.addEventListener("mousewheel", scrollHandler)
  // document.addEventListener("DOMMouseScroll", scrollHandler)
  // document.addEventListener("touchmove", scrollHandler)
});

  /**
  * Lazy load initial images in rendered markup
  * @param {Object} img - an array-like list of <img> DOM nodes in the initial index.html
  *
  * @return {Object}
  */
  function lazyLoad(img) {
    img.src = img.dataset.src;
    return img;
  }

/**
* Convenience function for building the api route
* @param {String} rover - the rover whose data we want
* @param {String} page - the page of data we want
*
* @return {String}
*/
function _mkRoute(rover, limit, page) {
  return `/rover/${rover}/limit/${limit}/page/${page}`
}

