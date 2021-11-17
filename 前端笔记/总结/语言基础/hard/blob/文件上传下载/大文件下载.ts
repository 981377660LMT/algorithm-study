// 1. 获取文件长度
// 通过发送 HEAD 请求，然后从响应头中读取 Content-Length 的信息，计算文件分块数
function getContentLength(url: string) {
  return new Promise<string | null>((resolve, reject) => {
    let xhr = new XMLHttpRequest()
    xhr.open('HEAD', url)
    xhr.send()
    xhr.onload = function () {
      resolve(xhr.getResponseHeader('Content-Length'))
    }
    xhr.onerror = reject
  })
}

// 2. 并发控制
// 实用asyncPool并发下载

// 3.下载指定范围内的数据块
function getBinaryContent(url: string, start: number, end: number, i: number) {
  return new Promise((resolve, reject) => {
    try {
      const xhr = new XMLHttpRequest()
      xhr.open('GET', url, true)
      xhr.setRequestHeader('range', `bytes=${start}-${end}`) // 请求头上设置范围请求信息
      xhr.responseType = 'arraybuffer' // 设置返回的类型为arraybuffer
      xhr.onload = function () {
        resolve({
          index: i, // 文件块的索引
          buffer: xhr.response, // 范围请求对应的数据
        })
      }
      xhr.send()
    } catch (err: any) {
      reject(new Error(err))
    }
  })
}

// 4.合并Uint8Array
function concatenate(arrays: Uint8Array[]) {
  if (!arrays.length) return null
  let totalLength = arrays.reduce((acc, value) => acc + value.length, 0)
  let result = new Uint8Array(totalLength)
  let length = 0
  for (let array of arrays) {
    result.set(array, length)
    length += array.length
  }
  return result
}

// 5.实现客户端文件保存的功能
function saveAs({ name, buffers, mime = 'application/octet-stream' }) {
  const blob = new Blob([buffers], { type: mime })
  const blobUrl = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.download = name || Math.random()
  a.href = blobUrl
  a.click()
  URL.revokeObjectURL(blobUrl)
}
