import axios from 'axios'

// 通过a标签来进行下载
function download() {
  const image = document.getElementById('test-image') as HTMLImageElement
  const link = document.createElement('a')

  link.download = 'test.png'
  link.rel = 'noopener'

  // 跨域
  if (link.origin !== location.origin) {
    axios
      .get<Blob>('url', { responseType: 'blob' })
      .then(data => {
        // 内存中的大文件
        link.href = URL.createObjectURL(data.data)
        setTimeout(() => {
          link.dispatchEvent(new MouseEvent('click'))
        })
        setTimeout(() => {
          URL.revokeObjectURL(link.href)
        }, 40000)
      })
      .catch(e => {
        // 错误时在新窗口显示图片
        console.log(e)
        link.target = '_black'
        link.href = 'url'
        link.dispatchEvent(new MouseEvent('click'))
      })
  } else {
    // 同源
    link.href = image.src
    link.dispatchEvent(new MouseEvent('click'))
  }
}
// rel 表示“关系 (relationship) ”
