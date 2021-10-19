const smoothScroll = (element: string) =>
  document.querySelector(element)!.scrollIntoView({
    behavior: 'smooth',
  })
smoothScroll('#fooBar') // scrolls smoothly to the element with the id fooBar
smoothScroll('.fooBar')
// scrolls smoothly to the first element with a class of fooBar

// auto 的滚动是直接跳转 很生硬
// smooth 带动画
