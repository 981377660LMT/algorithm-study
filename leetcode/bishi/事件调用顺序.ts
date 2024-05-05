// 腾讯-同步任务异步任务调用顺序
/* eslint-disable @typescript-eslint/no-non-null-assertion */

import assert from 'assert'

// 前端程序中异步调用很常见，问题排查中需要分析方法调用及结束顺序。
// 现指定方法名为单字母( a-zA-Z)，小写字母表示同步方法，
// 大写字母表示异步方法，比如异步下载文件再调用回调。
// 可以简单认为异步方法在被调用时开始，
// 异步方法中的其他方法调用在外层方法结束了才开始(不管外层方法是同步还是异步)。
// !简单起见，一组方法调用中最多只有一个异步方法。
// a.bc : b和c在a方法中依次同步调用，a开始>b开始并结束>c开始并结束>a结束。
// a.Bc B.d∶B是异步调用，a开始>B开始>c开始并结束>a结束>d开始并结束>B结束

// !1.模拟事件循环
// 先执行所有的同步代码，然后循环清空异步任务队列(asyncQueue)，直到队列为空。
// !2.执行任务(dfs)时，如果碰到同步任务，就直接执行，
// 如果碰到异步任务，就将其放入异步任务队列。
function solve(input: string): string {
  const words = input.split(/\s+/)
  const adjMap = new Map<string, string[]>()
  words.forEach(word => {
    const [parent, children] = word.split('.')
    for (let i = 0; i < children.length; i++) {
      const child = children[i]
      !adjMap.has(parent) && adjMap.set(parent, [])
      adjMap.get(parent)!.push(child)
    }
  })

  const res: string[] = []
  const asyncQueue: string[] = []
  dfs(input[0], true) // !1.执行同步
  // !2.循环清空异步任务队列
  while (asyncQueue.length) {
    const cur = asyncQueue.shift()!
    dfs(cur, false)
  }

  return res.join('')

  // 执行某个任务
  function dfs(cur: string, isFirst: boolean): void {
    isFirst && res.push(cur)
    const children = adjMap.get(cur) || []
    for (const child of children) {
      if (isLowerCase(child)) {
        dfs(child, true) // !1.碰到同步任务，直接执行
      } else {
        res.push(child)
        asyncQueue.push(child) // !2.碰到异步任务，放入异步任务队列
      }
    }
    res.push(cur)
  }
}

function isLowerCase(str: string): boolean {
  return str === str.toLowerCase()
}

assert.strictEqual(solve('a.bc  b.d'), 'abddbcca')
assert.strictEqual(solve('a.Bc  B.d'), 'aBccaddB')
assert.strictEqual(solve('A.B B.C C.D D.E'), 'ABACBDCEDE')

export {}
