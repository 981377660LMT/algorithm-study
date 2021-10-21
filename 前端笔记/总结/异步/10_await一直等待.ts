async function async1() {
  console.log('async1 start') // 2

  await new Promise(resolve => {
    console.log('promise1') // 3
  })

  console.log('async1 success')

  return 'async1 end'
}

console.log('srcipt start') // 1

async1().then(res => console.log(res))

console.log('srcipt end') // 4

// 在async1中await后面的Promise是没有返回值的，
// 也就是它的状态始终是pending状态，因此相当于一直在await，await，await却始终没有响应...
