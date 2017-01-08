/**
* Fetch images and metadata and append to the DOM
* @param {string} uri - a generator function to run asynchronously
*/
export function* update(uri) {
  if (~State.fetched.indexOf(uri)) {
    try {
      let res = yield fetch(uri) // returns a promise for the response
      let list = yield res.json() // returns a promise for json
      // move to next rover if this one is out of images
      if (!list || list.length < 10) {
        State.rover++
        if (State.rover > ROVERS.length) {
          State.rover = 0
        }
        State.page = 0
        return
      }
      for (let item of list) {
        let node = mkNode(item)
        State.nodes.push(node)
        State.main.append(node)
      }
      State.fetched.push(uri)
      State.page++
    } catch (err) {
      console.error(err)
    }
  }
}

/**
* Convenience function for building the api route
* @param {String} rover - the rover whose data we want
* @param {String} page - the page of data we want
*
* @return {String}
*/
export function mkRoute (rover, page) {
  return `/rover/${rover}/page/${page}`
}

/**
* Make new dom nodes with image data
* @param {Object} data - data for one image container div
*
* @return {Object}
*/
function mkNode(data) {
  let div = document.createElement('div')
  div.classList.add('wrapper-item', 'hidden')
  const tmpl = p => template`
  <div class="img-container">
    <img src="${p.img_src}" class="photo"></img>
  </div>
  <div class="metadata">
    <ul>
      <li>image id: ${p.id}</li>
      <li>martian sol: ${p.sol}</li>
      <li>earth date: ${p.earth_date}</li>
      <li>rover: ${p.rover}</li>
      <li>camera: ${p.camera}</li>
      <li><button class="ctrl-btn" type="button" name="backward" id="ctrl-backward">regress</button></li>
      <li><button class="ctrl-btn" type="button" name="forward" id="ctrl-forward">advance</button></li>
    </ul>
  </div>
  `
  div.innerHTML = tmpl(data)
  div.dataset.id = data.id
  return div
}
