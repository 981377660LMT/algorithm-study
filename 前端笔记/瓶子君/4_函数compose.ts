// 扁平化
const flattenDeep = (array: number[]) => array.flat(Infinity)

// 去重
const unique = (array: number[]) => [...new Set(array)]

// 排序
const sort = (array: number[]) => array.sort((a, b) => a - b)

const flatten_unique_sort = compose(flattenDeep, unique, sort)

// 测试
const arr = [[1, 2, 2], [3, 4, 5, 5], [6, 7, 8, 9, [11, 12, [12, 13, [14]]]], 10]
console.log(flatten_unique_sort(arr))

function compose(...funcs: Function[]) {
  const composeTwo = (func1: Function, func2: Function) => {
    return (...args: any[]) => func2(func1(...args))
  }

  return funcs.reduce(composeTwo)
}

export {}
