// start lazy loading the first images right away
let photos = Array.from(document.querySelectorAll('img.photo'), (img) => lazyLoad(img))

/**
* Lazy load initial images in rendered markup
* @param {Object} img - an array-like list of <img> DOM nodes in the initial index.html
*
* @return {Object}
*/
function lazyLoad(img) {
  img.src = img.dataset.src
  return img
}
