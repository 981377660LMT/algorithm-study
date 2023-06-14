// 现给出一个N行M列矩阵(N,M<=2e5)，并给出Q次询问，如下:
// 1.对列给[left,right]所有列元素加上x  形如 1 left right x
// 2.对行将所有第i行元素赋值为x  形如 2 i x
// 3.输出第i行j列元素  形如 3 i j

// AtCoder Beginner Contest 253 - SGColin的文章 - 知乎
// https://zhuanlan.zhihu.com/p/521516671

// !离线的做法就是把后面的贡献写做前缀和差分，然后两个时刻维护一下。
// !在线的做法就是写一个主席树 + 标记持久化；
// 记录每行最后一次被赋值的时间戳timei。和赋值vali
// 则答案为vali+[timei,now]这段操作里对yi加的值。
import * as fs from 'fs'

function useInput(debugCase?: string) {
  const data = debugCase == void 0 ? fs.readFileSync(process.stdin.fd, 'utf8') : debugCase
  const dataIter = _makeIter(data)

  function input(): string {
    return dataIter.next().value.trim()
  }

  function* _makeIter(str: string): Generator<string, string, undefined> {
    yield* str.trim().split(/\r\n|\r|\n/)
    return ''
  }

  return {
    input
  }
}

const { input } = useInput()
const [ROW, COL, q] = input().split(' ').map(Number)
const last = new Uint32Array(ROW + 10) // 记录每行最后一次被赋值的时间

for (let i = 1; i < q + 1; i++) {
  const [op, ...args] = input().split(' ').map(Number)
  if (op === 1) {
    const [left, right, delta] = args
  } else if (op === 2) {
    const [row, target] = args
  } else {
    const [row, col] = args
  }
}

// TODO
