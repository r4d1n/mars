'use strict';

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
  * HTML Templating
  *
  * A quick and dirty ES6 solution via http://www.2ality.com/2015/01/template-strings-html.html
  *
  * @param  {Function} literalSections - ES6 template literal strings for the template
  * @param  {Array} wait - time between function calls
  *
  * @return {Function}
  */
  function html(literalSections, ...substs) {
    // Use raw literal sections: we don’t want
    // backslashes (\n etc.) to be interpreted
    let raw = literalSections.raw

    let result = ''

    substs.forEach((subst, i) => {
      // Retrieve the literal section preceding
      // the current substitution
      let lit = raw[i]

      // In the example, map() returns an array:
      // If substitution is an array (and not a string),
      // we turn it into a string
      if (Array.isArray(subst)) {
        subst = subst.join('')
      }

      // If the substitution is preceded by a dollar sign,
      // we escape special characters in it
      if (lit.endsWith('$')) {
        subst = htmlEscape(subst)
        lit = lit.slice(0, -1)
      }
      result += lit
      result += subst
    })
    // Take care of last literal section
    // (Never fails, because an empty template string
    // produces one literal section, an empty string)
    result += raw[raw.length-1] // (A)

    return result
  }

  // public API
  return {
    generator: run,
    template: html
  }
})()
