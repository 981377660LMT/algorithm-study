// // JavaScript(V8 6.0.0)——更简单一点
// const json = JSON.parse(readline())
// let depth = Number(readline())

// // const json = [1, [2, [3, [4]], 5]]
// // let depth = 2
// let queue = json

// while (queue.length && depth--) {
//   const nextQueue = []

//   const length = queue.length
//   for (let _ = 0; _ < length; _++) {
//     const cur = queue.shift()
//     if (Array.isArray(cur)) {
//       nextQueue.push(...cur)
//     } else {
//       nextQueue.push(cur)
//     }
//   }

//   queue = nextQueue
// }

// // console.log(queue)
// print(JSON.stringify(queue))
////////////////////////////////////////
// let list = []
// let nums = []
// let [cur, pre, next] = ['', '', '']

// while ((line = readline())) {
//   const lines = line.split(',')
//   if (lines.length === 2) {
//     list.push(lines[0])
//     nums.push(Number(lines[1]))
//   } else {
//     cur = lines[0]
//     pre = lines[1]
//     next = lines[2]
//   }
// }

// // list = ['a', 'b', 'c', 'd', 'e', 'f']
// // ;[cur, pre, next] = ['a', 'f', '0']
// // nums = [1, 3, 6, 8, 9, 19]

// // list = ['a', 'b', 'c', 'd', 'e', 'f']
// // ;[cur, pre, next] = ['e', 'a', 'b']
// // nums = [1, 3, 6, 8, 9, 19]

// const rawIndex = list.indexOf(cur)
// const deleted = list.splice(rawIndex, 1)[0]
// const targetIndex = pre === '0' ? 0 : list.indexOf(pre) + 1
// // console.log(rawIndex, targetIndex, list)
// list.splice(targetIndex, 0, deleted)
// // console.log(list)

// for (let index = 0; index < list.length; index++) {
//   const element = list[index]
//   const sort = nums[index]
//   console.log(`${element},${sort}`)
// }
////////////////////////////////////////
