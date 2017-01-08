// immediately lazy load the first image
document.addEventListener('DOMContentLoaded', function() {
  lazyLoad(document.querySelector('img.photo'))
})

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
