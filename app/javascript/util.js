"use strict";

/**
* A helper for doing async tasks with generators
* @param {function} generatorFn - a generator function to run asynchronously
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

module.exports =  {
  runGenerator: run
}
