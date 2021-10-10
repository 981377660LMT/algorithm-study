type Command = '$push' | '$set' | '$merge' | '$apply'

/**
 * @param {any} data
 * @param {Object} command
 */
function update(data: any, command: Record<PropertyKey, any>) {
  for (const key of Object.keys(command)) {
    switch (key) {
      case '$push':
        return [...data, ...command[key]]
      case '$set':
        return command[key]
      case '$merge':
        if (typeof data !== 'object') throw new Error('invalid input')
        return { ...data, ...command[key] }
      case '$apply':
        return command[key](data)
      // 对象与数组的key
      default:
        data[key] = update(data[key], command[key])
    }
  }

  return data
}

if (require.main === module) {
  // 1. {$push: array} 添加元素到数组
  const arr = [1, 2, 3, 4]
  const newArr = update(arr, { $push: [5, 6] })
  // // [1, 2, 3, 4, 5, 6]

  // 2. {$set: any} 修改目标
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
  // /*
  // {
  //   a: {
  //     b: {
  //       c: 3
  //     }
  //   },
  //   d: 2
  // }
  // */

  // 3. {$merge: object} 合并到目标object
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
  // /*
  // {
  //   a: {
  //     b: {
  //       c: 1,
  //       e: 5
  //     }
  //   },
  //   d: 2
  // }
  // */

  // 4. {$apply: function} 自定义修改
  // const arr = [1, 2, 3, 4]
  // const newArr = update(arr, {0: {$apply: (item) => item * 2}})
  // // [2, 2, 3, 4]
}

export {}
