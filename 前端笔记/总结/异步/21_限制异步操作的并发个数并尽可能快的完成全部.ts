/* eslint-disable func-names */
/* eslint-disable no-await-in-loop */

type AsyncFunction<T = unknown> = (...args: unknown[]) => Promise<T>

function throttlePromises<T>(funs: AsyncFunction<T>[], max = 3): Promise<T[]> {
  const res: T[] = Array(funs.length).fill(null)
  const iter = funs.entries()

  // 注意他们共享一个iter 而iter是会被耗尽的 所以他们执行的是不同的promise
  // 即：使用fill来共用iter的引用
  // !Promise.all 并行三个槽 但是在每个槽里的任务是串行的
  const workers = Array(max).fill(iter).map(work)
  return Promise.all(workers).then(() => res)

  async function work(entries: IterableIterator<[number, AsyncFunction<T>]>) {
    for (const [index, promiseFunc] of entries) {
      const cur = await promiseFunc()
      res[index] = cur
    }
  }
}

function loadImg(url: string) {
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.src = url

    img.onload = function () {
      console.log('加载完成')
      resolve(img)
    }

    img.onerror = function () {
      reject(new Error('加载失败'))
    }
  })
}

const urls = [
  'https://hexo-blog-1256114407.cos.ap-shenzhen-fsi.myqcloud.com/AboutMe-painting1.png',
  'https://hexo-blog-1256114407.cos.ap-shenzhen-fsi.myqcloud.com/AboutMe-painting2.png',
  'https://hexo-blog-1256114407.cos.ap-shenzhen-fsi.myqcloud.com/AboutMe-painting3.png',
  'https://hexo-blog-1256114407.cos.ap-shenzhen-fsi.myqcloud.com/AboutMe-painting4.png',
  'https://hexo-blog-1256114407.cos.ap-shenzhen-fsi.myqcloud.com/AboutMe-painting5.png',
  'https://hexo-blog-1256114407.cos.ap-shenzhen-fsi.myqcloud.com/bpmn6.png',
  'https://hexo-blog-1256114407.cos.ap-shenzhen-fsi.myqcloud.com/bpmn7.png',
  'https://hexo-blog-1256114407.cos.ap-shenzhen-fsi.myqcloud.com/bpmn8.png'
]

throttlePromises(urls.map(url => () => loadImg(url)))

export {}
