// https://zhuanlan.zhihu.com/p/49077183

// 解构速度测试(desctructing)
// !1.不用解构是最快的
// !2.对象解构比数组解构快很多(后者需要调用迭代器)

const arr1 = Array(1e7)
for (let i = 0; i < arr1.length; i++) arr1[i] = [i, i + 1]

const time1 = performance.now()
for (let i = 0; i < arr1.length; i++) {
  const [a, b] = arr1[i]
  a, b
}
console.log(performance.now() - time1) // 28.199999999254942

const time2 = performance.now()
for (let i = 0; i < arr1.length; i++) {
  const { 0: a, 1: b } = arr1[i]
  a, b
}
console.log(performance.now() - time2) // 24.5

const time3 = performance.now()
for (let i = 0; i < arr1.length; i++) {
  const item = arr1[i]
  const a = item[0]
  const b = item[1]
  a, b
}
console.log(performance.now() - time3) // 21

const time4 = performance.now()
for (let i = 0; i < arr1.length; i++) {
  const a = arr1[i][0]
  const b = arr1[i][1]
  a, b
}
console.log(performance.now() - time4) // 20.699999999254942

const time5 = performance.now()
arr1.forEach(item => {
  const a = item[0]
  const b = item[1]
  a, b
})
console.log(performance.now() - time5) // 93.79999999701977
