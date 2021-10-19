// 使用 Element.getBoundingClientRect () 和
// Window.inner (Width | Height)值来确定给定的元素在 viewport 中是否可见。

declare const el: Element

elementIsVisibleInViewport(el) // false - (not fully visible)

function elementIsVisibleInViewport(el: Element, partiallyVisible = false) {
  const { top, left, bottom, right } = el.getBoundingClientRect()
  const { innerHeight, innerWidth } = window

  return partiallyVisible
    ? ((top > 0 && top < innerHeight) || (bottom > 0 && bottom < innerHeight)) &&
        ((left > 0 && left < innerWidth) || (right > 0 && right < innerWidth))
    : top >= 0 && left >= 0 && bottom <= innerHeight && right <= innerWidth
}

export {}
