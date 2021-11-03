// 渲染百万条结构简单的大数据时 怎么使用分片思想优化渲染
// 扩展思考：对于大数据量的简单 dom 结构渲染可以用分片思想解决 如果是复杂的 dom 结构渲染如何处理？
// 这时候就需要使用虚拟列表和虚拟表格了
const ul = document.getElementById('root')!
const total = 100000
const once = 20

function render(cur: number, remain: number) {
  if (remain <= 0) return

  const pageCount = Math.min(remain, once)

  const frameRequestCallback = () => {
    for (let i = 0; i < pageCount; i++) {
      const li = document.createElement('li')
      li.innerText = cur + i + ' : ' + ~~(Math.random() * total)
      ul.appendChild(li)
    }

    render(cur + pageCount, remain - pageCount)
  }

  window.requestAnimationFrame(frameRequestCallback)
}

render(0, total)

export {}
