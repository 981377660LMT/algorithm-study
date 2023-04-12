function getUrlParam(sUrl: string, sKey: string) {
  return new URL(sUrl).searchParams.getAll(sKey)
}

// console.log(getUrlParam('http://www.nowcoder.com?key=1&key=2&key=3&test=4#hehe', 'key'))
/////////////////////////////////////////////////////////////////////////////
// 查找两个节点的最近的一个共同父节点，可以包括节点自身
function commonParentNode(oNode1: HTMLElement, oNode2: HTMLElement): HTMLElement {
  if (oNode1.contains(oNode2)) {
    return oNode1
  } else {
    return commonParentNode(oNode1.parentElement!, oNode2)
  }
}
/////////////////////////////////////////////////////////////////////////
// 根据包名，在指定空间中创建对象
// 输入描述：
// namespace({a: {test: 1, b: 2}}, 'a.b.c.d')
// 输出描述：
// {a: {test: 1, b: {c: {d: {}}}}}
function namespace(oNamespace: Record<string, any>, sPackage: string) {
  const chars = sPackage.split('.')
  let root = oNamespace

  for (const char of chars) {
    if (typeof root[char] !== 'object') root[char] = {}
    root = root[char]
  }

  return oNamespace
}
// console.dir(namespace({ a: { test: 1, b: 2 } }, 'a.b.c.d'), { depth: null })
/////////////////////////////////////////////////////////////////////////
// 数组去重
const uniq = function (arr: any[]) {
  return [...new Set(arr)]
}
// console.log(uniq([false, true, undefined, null, NaN, 0, 1, {}, {}, 'a', 'a', NaN]))
/////////////////////////////////////////////////////////////////////////
// 时间格式化输出
function formatDate(t: Date, str: string) {
  // console.log(t.getFullYear(), t.getMonth(), t.getHours(), t.getDay())
  const strategy = {
    yyyy: t.getFullYear(),
    yy: t.getFullYear().toString().slice(-2),
    M: t.getMonth() + 1,
    MM: (t.getMonth() + 1).toString().padStart(2, '0'),
    d: t.getDay(),
    dd: t.getDay().toString().padStart(2, '0'),
    H: t.getHours(),
    HH: t.getHours().toString().padStart(2, '0'),
    h: t.getHours() % 12,
    hh: (t.getHours() % 12).toString().padStart(2, '0'),
    m: t.getMinutes(),
    mm: t.getMinutes().toString().padStart(2, '0'),
    s: t.getSeconds(),
    ss: t.getSeconds().toString().padStart(2, '0'),
    w: ['日', '一', '二', '三', '四', '五', '六'][t.getDay()]
  } as Record<string, any>

  return str.replace(/(\w+)/gi, (_, group) => {
    return strategy[group]
  })
}
// console.log(formatDate(new Date(1409894060000), 'yyyy-MM-dd HH:mm:ss 星期w'))
// 输出：
// 2014-09-05 13:14:20 星期五
/////////////////////////////////////////////////////////////////////////
// 邮箱字符串判断
function isAvailableEmail(sEmail: string) {
  return /^([\w+\.])+@(\w+)([.]\w+)+$/.test(sEmail)
}
// console.log(isAvailableEmail('test@qq.com'))
/////////////////////////////////////////////////////////////////////////
// 将字符串转换为驼峰格式
function cssStyle2DomStyle(sName: string) {
  const reg = /-(\w)/g
  return sName.replace(reg, (match, g1, index) => {
    if (index === 0) return g1
    return g1.toUpperCase()
  })
}
// console.log(cssStyle2DomStyle('font-size'))
// console.log(cssStyle2DomStyle('-webkit-border-image'))
/////////////////////////////////////////////////////////////////////////
function removeWithoutCopy(arr: number[], item: number) {
  let l = 0

  for (let r = 0; r < arr.length; r++) {
    if (arr[r] === item) continue
    arr[l] = arr[r]
    l++
  }
  return arr.slice(0, l)
}
// console.log(removeWithoutCopy([1, 2, 2, 3, 4, 2, 2], 2))
/////////////////////////////////////////////////////////////////////////
// 打点计时器
// 1、从 start 到 end（包含 start 和 end），每隔 100 毫秒 console.log 一个数字，每次数字增幅为 1
// 2、返回的对象中需要包含一个 cancel 方法，用于停止定时操作
// 3、第一个数需要立即输出
// 通过setInterval()方法
function count(start: number, end: number) {
  console.log(start++)

  const timer = setInterval(() => {
    if (start <= end) console.log(start++)
    else clearInterval(timer)
  })

  const cancel = () => clearInterval(timer)

  return { cancel }
}

// 通过setTimeout()方法
function count2(start: number, end: number) {
  let timer: ReturnType<typeof setTimeout> | null

  if (start <= end) {
    console.log(start++)
    timer = setTimeout(() => {
      count2(start, end)
    }, 100)
  }

  const cancel = () => clearInterval(timer)

  return { cancel }
}

// const timer = count2(0, 10)
// timer.cancel()

/////////////////////////////////////////////////////////////////////////
function functionFunction(str: string) {
  return (...args: any[]) => `${str},${args.join(',')}`
}
// console.log(functionFunction('hello')('world', 'kk'))
/////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////

export {}
