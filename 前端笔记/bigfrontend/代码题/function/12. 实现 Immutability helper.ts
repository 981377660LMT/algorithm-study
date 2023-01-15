/* eslint-disable no-unreachable-loop */
// 12. 实现 Immutability helper

type Command = {
  $push?: unknown
  $set?: unknown
  $merge?: object
  $apply?: (arg: unknown) => unknown
  [key: string]: unknown
}

// dfs and mutate data
function update(data: Record<PropertyKey, unknown>, command: Command): unknown {
  for (const key of Object.keys(command)) {
    switch (key) {
      case '$push':
        return [...data, ...command[key]]
      case '$set':
        return command[key]
      case '$merge':
        if (typeof data !== 'object') throw new Error('data is not an object')
        return { ...data, ...command[key] }
      case '$apply':
        return command[key](data)
      default:
        data[key] = update(data[key], command[key])
    }
  }
  return data
}

// 如果你使用React，你肯定会遇到想要修改state的一部分的情况。
// 比如下面的state。
// const state = {
//   a: {
//     b: {
//       c: 1
//     }
//   },
//   d: 2
// }
// 如果我们想要修改d来生成一个新的state，我们可以用 _.cloneDeep，
// 但是这样没必要因为state.a并不需要被clone。
// 一个更好的办法是如下的浅拷贝
// const newState = {
//   ...state,
//   d: 3
// }
// 但是又有了新问题，如果我们同时需要修改c的话，我们需要写很复杂的代码，比如：
// const newState = {
//   ...state,
//   a: {
//     ...state.a,
//     b: {
//        ...state.b,
//        c: 2
//     }
//   }
// }
// 这显然还不如cloneDeep。
// Immutability Helper 可以很好的解决这个问题。
// 请实现你自己的Immutability helper update()，需要支持如下调用
// !1. {$push: array} 添加元素到数组
// const arr = [1, 2, 3, 4]
// const newArr = update(arr, {$push: [5, 6]})
// [1, 2, 3, 4, 5, 6]
// !2. {$set: any} 修改目标
// const state = {
//   a: {
//     b: {
//       c: 1
//     }
//   },
//   d: 2
// }

// const newState = update(
//   state,
//   {a: {b: { c: {$set: 3}}}}
// )
//
// {
//   a: {
//     b: {
//       c: 3
//     }
//   },
//   d: 2
// }
//
// !3. {$merge: object} 合并到目标object
// const state = {
//   a: {
//     b: {
//       c: 1
//     }
//   },
//   d: 2
// }

// const newState = update(
//   state,
//   {a: {b: { $merge: {e: 5}}}}
// )
//
// {
//   a: {
//     b: {
//       c: 1,
//       e: 5
//     }
//   },
//   d: 2
// }
//
// !4. {$apply: function} 自定义修改
// const arr = [1, 2, 3, 4]
// const newArr = update(arr, {0: {$apply: (item) => item * 2}})
// [2, 2, 3, 4]
