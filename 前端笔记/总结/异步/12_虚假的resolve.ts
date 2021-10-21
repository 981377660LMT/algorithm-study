async function async1() {
  console.log('async1 start') // 2
  await new Promise(resolve => {
    console.log('promise1') // 3
    resolve('promise resolve')
  })
  console.log('async1 success') // 5
  return 'async1 end'
}
console.log('srcipt start') // 1
async1().then(res => {
  console.log(res) // 6  'async1 end'
})
new Promise(resolve => {
  console.log('promise2') // 4
  setTimeout(() => {
    console.log('timer') // 7
  })
})
// 在async1中的new Promise它的resovle的值和async1().then()里的值是没有关系的，
// 很多小伙伴可能看到resovle('promise resolve')就会误以为是async1().then()中的返回值。
