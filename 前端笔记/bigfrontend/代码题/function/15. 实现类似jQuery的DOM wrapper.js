/**
 * @param {HTMLElement} el - element to be wrapped
 */
function $(el) {
  // your code herer
  return new Wrapper(el)
}

class Wrapper {
  constructor(el) {
    this.el = el
  }

  css(prop, value) {
    this.el.style[prop] = value
    return this
  }
}

if (require.main === module) {
  $('#button').css('color', '#fff').css('backgroundColor', '#000').css('fontWeight', 'bold')
}
