// 请实现一个类似于 jQuery.on()的方法，用来注册事件处理器。
// 在jQuery中，可以用selector来选择element，在本题目中，这简化为一个predicate函数。

/**
 * @param {HTMLElement} root
 * @param {(el: HTMLElement) => boolean} predicate
 * @param {(e: Event) => void} handler
 */
function onClick(
  root: HTMLElement,
  predicate: (el: HTMLElement) => boolean,
  handler: (e: Event) => void
) {
  // your code here
  root.hand
}

if (require.main === module) {
  onClick(
    // root element
    document.body,
    // predicate
    el => el.tagName.toLowerCase() === 'div',
    function (e) {
      console.log(this)
      // this logs all the `div` element
    }
  )
}

// event.stopPropagation() 和 event.stopImmediatePropagation() 都需要被支持。
// event.stopPropagation阻止捕获和冒泡阶段中当前事件的进一步传播。
// Event 接口的 stopImmediatePropagation() 方法阻止监听同一事件的其他事件监听器被调用

// 请只在root element上添加一个真正的event listener
