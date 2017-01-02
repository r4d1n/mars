"use strict";

const Util = (function() {
  /**
  * A helper for doing async tasks with generators
  * @param {Function} generatorFn - a generator function to run asynchronously
  * Additional arguments are passed to the generator
  */
  let iterator

  function run(generatorFn) {
    let args = [].slice.call(arguments, 1)
    // initialize the generator in the current context with args
    iterator = generatorFn.apply(this, args)
    return Promise.resolve()
    .then(() => handleResult(iterator.next()))
  }

  function handleResult(next){
    if (next.done) {
      return next.value
    } else {
      return Promise.resolve(next.value)
      // pass current value back to the generator and recurse with what comes back
      .then((val) => {
        return handleResult(iterator.next(val))
      })
      .catch((err) => {
        // pass error back into the generator
        return iterator.throw(err)
      })
    }
  }

  /**
  * Lazy load initial images in rendered markup
  * @param {Object} img - an array-like list of <img> DOM nodes in the initial index.html
  */
  function lazyLoad(img) {
    img.src = img.dataset.src
    return img
  }

  /**
  * Convenience function for building the api route
  * @param {String} rover - the rover whose data we want
  * @param {String} page - the page of data we want
  */
  function makeRoute (rover, page) {
    return `/rover/${rover}/page/${page}`
  }

  /**
  * Debounce
  *
  * @param  {Function} fn the function to call
  * @param  {Number} wait - time between function calls
  * @param  {Boolean} immediate - trigger before the wait instead of after
  *
  * @return {Function}
  */
  function debounce(fn, wait, immediate) {
    let timeout = null

    return function () {
      let self = this
      let args = arguments

      clearTimeout(timeout)

      timeout = setTimeout(function () {
        timeout = null
        if (!immediate) {
          fn.apply(self, args)
        }
      }, wait)

      if (immediate && !timeout) {
        fn.apply(self, args)
      }
    };
  }

  // public API
  return {
    generator: run,
    lazyLoad: lazyLoad,
    makeRoute: makeRoute,
    debounce: debounce
  }
})()
