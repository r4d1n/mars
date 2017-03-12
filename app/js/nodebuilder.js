import "whatwg-fetch"

export class NodeBuilder {
  constructor() {}

  getData(uri) {
    fetch(uri)
      .then((res) => res.json())
      .catch((err) => console.error(err));
  }

  /**
  * Make and append a dom node for a new item
  * @param {object} data - the data to render
  * @param {Element} target - the DOM element to append to
  */
  render(data, target) {
    let node = this._mkNode(data);
    target.append(node);
  }

  /**
  * Make new dom nodes with image data
  * @param {Object} data - data for one image container div
  *
  * @return {Element}
  */
  _mkNode(data) {
    let div = document.createElement("div");
    div.classList.add("wrapper-item");
    const tmpl = p => this._template(`<div class="img-container">
        <img src="${p.img_src}" class="photo"></img>
      </div>
      <div class="metadata">
        <ul>
          <li>image id: ${p.id}</li>
          <li>martian sol: ${p.sol}</li>
          <li>earth date: ${p.earth_date}</li>
          <li>rover: ${p.rover}</li>
          <li>camera: ${p.camera}</li>
        </ul>
      </div>`);
    div.innerHTML = tmpl(data);
    div.dataset.id = data.id;
    return div;
  }

  /**
  * Make new dom nodes without data
  *
  * @return {Element}
  */
  createPlaceholder() {
    let div = document.createElement('div');
    div.classList.add('.wrapper-item');
    return div;
  }

  /**
    * HTML Templating
    *
    * A quick and dirty ES6 solution via http://www.2ality.com/2015/01/template-strings-html.html
    *
    * @param  {Function} literalSections - ES6 template literal strings for the template
    * @param  {Array} substs - the data that will fill in the template
    *
    * @return {Function}
    */
  _template(literalSections, ...substs) {
    // Use raw literal sections: we don’t want
    // backslashes (\n etc.) to be interpreted
    let raw = literalSections.raw;

    let result = '';

    substs.forEach((subst, i) => {
      // Retrieve the literal section preceding
      // the current substitution
      let lit = raw[i];

      // In the example, map() returns an array:
      // If substitution is an array (and not a string),
      // we turn it into a string
      if (Array.isArray(subst)) {
        subst = subst.join("");
      }

      // If the substitution is preceded by a dollar sign,
      // we escape special characters in it
      if (lit.endsWith("$")) {
        subst = htmlEscape(subst);
        lit = lit.slice(0, -1);
      }
      result += lit;
      result += subst;
    })
    // Take care of last literal section
    // (Never fails, because an empty template string
    // produces one literal section, an empty string)
    result += raw[raw.length - 1]; // (A)

    return result;
  }
}
