// 利用 XMLHttpRequest 手写 AJAX 实现
const get = (url: string) => {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()
    xhr.open('get', url, true)
    xhr.setRequestHeader('Content-type', 'application/json')
    xhr.onreadystatechange = () => {
      if (xhr.readyState !== 4) return
      if (xhr.status === 200 || xhr.status === 304) {
        resolve(xhr.responseText)
      } else {
        reject(new Error(xhr.responseText))
      }
    }

    xhr.send()
  })
}

// 只要 readyState 属性发生变化
// 就会调用XMLHttpRequest.onreadystatechange
