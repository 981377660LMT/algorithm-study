// 使用 Window.getComputedStyle ()获取指定元素的 CSS 规则值。
getStyle(document.querySelector('p')!, 'fontSize') // '16px'

function getStyle(el: HTMLElement, prop: keyof CSSStyleDeclaration) {
  return getComputedStyle(el)[prop]
}
