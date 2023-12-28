/* eslint-disable arrow-body-style */

// !每次从未处理的数组中随机取一个元素，然后把该元素放到数组的尾部，即数组的尾部放的就是已经处理过的元素

/**
 * 用 Fisher-Yates 方法随机打乱数组。
 */
function shuffle<T>(arr: T[]): void {
  const randInt = (min: number, max: number): number => {
    return min + Math.floor((max - min + 1) * Math.random())
  }
  const swap = (arr: T[], i: number, j: number): void => {
    const tmp = arr[i]
    arr[i] = arr[j]
    arr[j] = tmp
  }

  for (let i = 0; i < arr.length; i++) {
    const rand = randInt(0, i)
    swap(arr, i, rand)
  }
}

if (require.main === module) {
  const foo = [1, 2, 3, 4, 5, 6]
  shuffle(foo)
  console.log(foo)
}

export { shuffle }
