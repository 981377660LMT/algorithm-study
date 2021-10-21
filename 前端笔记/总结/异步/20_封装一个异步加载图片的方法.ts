type AsyncFunction<T = any> = (...args: any[]) => Promise<T>
const urls = [
  'https://hexo-blog-1256114407.cos.ap-shenzhen-fsi.myqcloud.com/AboutMe-painting1.png',
  'https://hexo-blog-1256114407.cos.ap-shenzhen-fsi.myqcloud.com/AboutMe-painting2.png',
  'https://hexo-blog-1256114407.cos.ap-shenzhen-fsi.myqcloud.com/AboutMe-painting3.png',
  'https://hexo-blog-1256114407.cos.ap-shenzhen-fsi.myqcloud.com/AboutMe-painting4.png',
  'https://hexo-blog-1256114407.cos.ap-shenzhen-fsi.myqcloud.com/AboutMe-painting5.png',
  'https://hexo-blog-1256114407.cos.ap-shenzhen-fsi.myqcloud.com/bpmn6.png',
  'https://hexo-blog-1256114407.cos.ap-shenzhen-fsi.myqcloud.com/bpmn7.png',
  'https://hexo-blog-1256114407.cos.ap-shenzhen-fsi.myqcloud.com/bpmn8.png',
]

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

function throttlePromises<T>(funcs: AsyncFunction<T>[], max = 3): Promise<T[]> {
  const results: T[] = []
  const iter = funcs.entries()

  // 注意他们共享一个iter 而iter是会被耗尽的 所以他们执行的是不同的promise
  // 即：使用fill来共用iter的引用
  const workers = Array(max).fill(iter).map(work)

  return Promise.all<void>(workers).then(() => results)

  async function work(entries: IterableIterator<[number, AsyncFunction<T>]>) {
    for (const [index, promiseFunc] of entries) {
      const result = await promiseFunc(1)
      results[index] = result
    }
  }
}

throttlePromises(urls.map(url => () => loadImg(url)))

export {}
