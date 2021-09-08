// Window 对象作为客户端 JavaScript 的全局对象，它有一个 window 属性指向自身
function isWindow(obj: any): obj is Window {
  return obj != null && obj === obj.window
}

function isElement(obj: Node): obj is Element {
  return !!(obj && obj.nodeType === 1)
}

export {}
