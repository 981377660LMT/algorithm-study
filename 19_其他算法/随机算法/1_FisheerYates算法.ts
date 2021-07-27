// 每次从未处理的数组中随机取一个元素，然后把该元素放到数组的尾部，即数组的尾部放的就是已经处理过的元素
// O(n)
// lodash的shuffle使用此方法
const shuffle = (arr: number[]): void => {
  let len = arr.length
  while (len) {
    const rand = Math.floor(Math.random() * len)
    ;[arr[len - 1], arr[rand]] = [arr[rand], arr[len - 1]]
    len--
  }
}

const foo = [1, 2, 3, 4, 5, 6]
shuffle(foo)
console.log(foo)

export {}
