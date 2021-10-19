// 顺滑:Window.requestAnimationFrame ()对滚动条进行动画处理。
scrollToTop() // Smooth-scrolls to the top of the page

function scrollToTop() {
  const diff = document.documentElement.scrollTop || document.body.scrollTop
  if (diff > 0) {
    window.requestAnimationFrame(scrollToTop) // 宏任务
    window.scrollTo(0, diff - diff / 8)
  }
}
