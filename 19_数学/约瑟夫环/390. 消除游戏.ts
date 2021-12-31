// 从左到右， 还是从右到左，每次都要消除 一半的数
// 关注何时消除队首元素:`从左到右` 或者 `从右到左，数组为奇数个`，才会消除第一个。

function lastRemaining(n: number): number {
  let remain = n
  let isToRight = true
  let res = 1
  let step = 1

  while (remain > 1) {
    if (isToRight || (remain & 1) === 1) {
      res += step
    }

    isToRight = !isToRight
    step *= 2
    remain >>>= 1
  }

  return res
}

console.log(lastRemaining(9))

// n = 9,
// 1 2 3 4 5 6 7 8 9
// 2 4 6 8
// 2 6
// 6

// 输出:
// 6
