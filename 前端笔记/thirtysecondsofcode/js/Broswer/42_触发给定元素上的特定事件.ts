triggerEvent(document.getElementById('myId')!, 'click')
triggerEvent(document.getElementById('myId')!, 'click', { username: 'bob' })

// 触发给定元素上的特定事件，可以选择传递自定义数据。
function triggerEvent(el: Element, eventType: string, detail?: any) {
  return el.dispatchEvent(new CustomEvent(eventType, { detail }))
}
