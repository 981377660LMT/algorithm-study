async function testSometing() {
  console.log('执行testSometing') // 2
  return 'testSometing'
}

async function testAsync() {
  console.log('执行testAsync') // 6
  return Promise.resolve('hello async')
}

async function test() {
  console.log('test start...') // 1
  const v1 = await testSometing()
  console.log(v1) // 5
  const v2 = await testAsync()
  console.log(v2) // 8
  console.log(v1, v2) // 9
}

test()

const promise = new Promise(resolve => {
  console.log('promise start...') // 3
  resolve('promise')
})
promise.then(val => console.log(val)) // 7

console.log('test end...') // 4

export {}
