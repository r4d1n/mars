import { NodeBuilder } from "./nodebuilder"

// Number of items to instantiate beyond current view in the scroll direction.
const RUNWAY_ITEMS = 50

// Number of items to instantiate beyond current view in the opposite direction.
const RUNWAY_ITEMS_REVERSE = 10

// The number of pixels of additional length to allow scrolling to.
const SCROLL_RUNWAY = 2000

// The animation interval (in ms) for fading in content from placeholders.
const ANIMATION_DURATION_MS = 200


export class InfiniteScroll {
  /**
   * Construct an infinite scroller.
   * @param {Element} element The scrollable element to use as the infinite scroll region.
   *
   */
  constructor(element) {
    this.anchorItem = { index: 0, offset: 0 }
    this._firstItem_ = 0
    this._lastItem = 0
    this.anchorScrollTop = 0
    this._region = element
    this._items = []
    this._loadedItems = 0
    this._requestInProgress = false
    this._region.addEventListener('scroll', this._onScroll.bind(this))
    // window.addEventListener('resize', this._onResize.bind(this))

    // Create an element to force the scroller to allow scrolling to a certain
    // point.
    this._runway = document.createElement('div')
    this._runwayLimit = 0
    this._runway.style.position = 'absolute'
    this._runway.style.height = '1px'
    this._runway.style.width = '1px'
    this._runway.style.transition = 'transform 0.2s'
    this._region.appendChild(this._runway)
    // this._onResize()
    console.log("infinite scroll", this)
  }

  /**
    * Scroll handler:
    * Determine the newly anchored item and offset
    * Update the visible elements
    * Requesting more items if necessary
    */
  _onScroll() {
    let delta = this._region.scrollTop - this.anchorScrollTop
    // Special case, if we get to very top, always scroll to top.
    if (this._region.scrollTop == 0) {
      this.anchorItem = { index: 0, offset: 0 }
    } else {
      this.anchorItem = this.calculateAnchoredItem(this.anchorItem, delta)
    }
    this.anchorScrollTop = this._region.scrollTop
    let lastItemOnScreen = this.calculateAnchoredItem(this.anchorItem, this._region.offsetHeight)
    if (delta < 0)
      this.fill(this.anchorItem.index - RUNWAY_ITEMS, lastItemOnScreen.index + RUNWAY_ITEMS_REVERSE)
    else
      this.fill(this.anchorItem.index - RUNWAY_ITEMS_REVERSE, lastItemOnScreen.index + RUNWAY_ITEMS)
  }
}
