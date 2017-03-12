import { NodeBuilder } from "./nodebuilder"

// Number of items to instantiate beyond current view in the scroll direction.
const RUNWAY_ITEMS = 30

// Number of items to instantiate beyond current view in the opposite direction.
const RUNWAY__itemsREVERSE = 20

// The number of pixels of additional length to allow scrolling to.
const SCROLL_RUNWAY = 2000

export class InfiniteScroll {
  /**
   * Construct an infinite scroller.
   * @param {Element} element The scrollable element to use as the infinite scroll region.
   *
   */
  constructor(element) {
    this.nodeBuilder = new NodeBuilder();
    this._region = element;
    this.anchorItem = { index: 0, offset: 0 };
    this._firstItem_ = 0;
    this._lastItem = 0;
    this.anchorScrollTop = 0;
    this._items = [];
    this._placeholders = [];
    this._loadedItems = 0;
    this._requestInProgress = false;
    this._region.addEventListener('scroll', this._onScroll.bind(this));
    window.addEventListener('resize', this._onResize.bind(this));

    this._runway = document.createElement('div');
    this._runwayLimit = 0;
    this._runway.style.position = 'absolute';
    this._runway.style.height = '1px';
    this._runway.style.width = '1px';
    this._runway.style.transition = 'transform 0.2s';
    this._region.appendChild(this._runway);
    this._onResize();
    console.log("infinite scroll", this);
  }

  /**
    * Scroll handler:
    * Determine the newly anchored item and offset
    * Update the visible elements
    * Request more items if necessary
    */
  _onScroll() {
    console.log('scroll fires')
    let delta = this._region.scrollTop - this.anchorScrollTop;
    console.log('this._region.scrollTop', this._region.scrollTop, 'this.anchorScrollTop:', this.anchorScrollTop, 'delta:', delta);
    // Special case, if we get to very top, always scroll to top.
    if (this._region.scrollTop == 0) {
      this.anchorItem = { index: 0, offset: 0 };
    } else {
      this.anchorItem = this.calculateAnchoredItem(this.anchorItem, delta);
    }
    this.anchorScrollTop = this._region.scrollTop;
    let lastItemOnScreen = this.calculateAnchoredItem(this.anchorItem, this._region.offsetHeight);
    console.log('lastItemOnScreen:', lastItemOnScreen);
    if (delta < 0)
      this.fill(this.anchorItem.index - RUNWAY_ITEMS, lastItemOnScreen.index + RUNWAY__itemsREVERSE);
    else
      this.fill(this.anchorItem.index - RUNWAY__itemsREVERSE, lastItemOnScreen.index + RUNWAY_ITEMS);
  }

  _onResize() {
    let node = this.nodeBuilder.createPlaceholder();
    node.style.position = 'absolute';
    this._region.appendChild(node);
    node.classList.remove('invisible');
    this._nodeSize = node.offsetHeight;
    this._nodeWidth = node.offsetWidth;
    this._region.removeChild(node);

    // Reset the cached sizes of items in the scroller
    this._items.forEach((el) => {
      el.height = 0;
      el.width = 0;
    });
    this._onScroll();
  }

  /**
   * Sets the range of items which should be attached and attaches those items.
   * @param {number} start The first item which should be attached.
   * @param {number} end One past the last item which should be attached.
   */
  fill(start, end) {
    this.firstAttachedItem_ = Math.max(0, start);
    this.lastAttachedItem_ = end;
    this.attachContent();
  }

  /**
   * Calculates the item that should be anchored after scrolling by delta from
   * the initial anchored item.
   * @param {{index: number, offset: number}} initialAnchor The initial position
   *     to scroll from before calculating the new anchor position.
   * @param {number} delta The offset from the initial item to scroll by.
   * @return {{index: number, offset: number}} Returns the new item and offset
   *     scroll should be anchored to.
   */
  calculateAnchoredItem(initialAnchor, delta) {
    if (delta == 0) return initialAnchor;
    delta += initialAnchor.offset;
    var i = initialAnchor.index;
    var tombstones = 0;
    if (delta < 0) {
      while (delta < 0 && i > 0 && this._items[i - 1].height) {
        delta += this._items[i - 1].height;
        i--;
      }
      tombstones = Math.max(-i, Math.ceil(Math.min(delta, 0) / this._nodeSize));
    } else {
      while (delta > 0 && i < this._items.length && this._items[i].height && this._items[i].height < delta) {
        delta -= this._items[i].height;
        i++;
      }
      if (i >= this._items.length || !this._items[i].height) {
        tombstones = Math.floor(Math.max(delta, 0) / this._nodeSize);
      }
    }
    i += tombstones;
    delta -= tombstones * this._nodeSize
    return {
      index: i,
      offset: delta
    }
  }

  /**
     * Creates or returns an existing tombstone ready to be reused.
     * @return {Element} A tombstone element ready to be used.
     */
  getPlaceholder() {
    var tombstone = this._placeholders.pop();
    if (tombstone) {
      tombstone.classList.remove('invisible');
      tombstone.style.opacity = 1;
      tombstone.style.transform = '';
      tombstone.style.transition = '';
      return tombstone;
    }
    return this.nodeBuilder.createPlaceholder();
  }

  /**
   * Attaches content to the scroller and updates the scroll position if
   * necessary.
   */
  attachContent() {
    console.log('attach content')
  }

  /**
   * Requests additional content if we don't have enough currently.
   */
  maybeRequestContent() {
    // Don't issue another request if one is already in progress as we don't
    // know where to start the next request yet.
    if (this._requestInProgress)
      return;
    var itemsNeeded = this.lastAttachedItem_ - this.loaded_items;
    if (itemsNeeded <= 0)
      return;
    this._requestInProgress = true;
    var lastItem = this._items[this.loaded_items - 1];
    this.source_.fetch(itemsNeeded).then(this.addContent.bind(this));
  }

  /**
   * Adds an item to the items list.
   */
  _addItem() {
    this._items.push({
      'data': null,
      'node': null,
      'height': 0,
      'width': 0,
      'top': 0,
    })
  }

  /**
   * Adds the given array of items to the items list and then calls
   * attachContent to update the displayed content.
   * @param {Array<Object>} items The array of items to be added to the infinite
   *     scroller list.
   */
  addContent(items) {
    this._requestInProgress = false;
    var startIndex = this._items.length;
    for (var i = 0; i < items.length; i++) {
      if (this._items.length <= this.loaded_items)
        this._addItem();
      this._items[this.loaded_items++].data = items[i];
    }
    this.attachContent();
  }
}

