// 使用 Element.matches ()检查任何给定的元素节点是否与提供的选择器匹配。
findClosestMatchingNode(document.querySelector('span'), 'body') // body

function findClosestMatchingNode(el: Element | null, target: string) {
  if (!el) return null

  while (el) {
    if (el.matches && el.matches(target)) return el
    el = el.parentElement
  }

  return null
}
