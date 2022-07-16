window.addEventListener('click', e => console.log(e), {
  capture: false,
  // passive如果被指定为true表示永远不会执行preventDefault(),
  // 这在实现丝滑柔顺的滚动的效果中很重要
  passive: false,
  once: false,
})
