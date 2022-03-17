function foo(this: any) {
  console.log(this.name)
  console.log(11)
}

// var obj = {
//   name: 'Heternally',
// }

// var obj1 = {
//   name: 'Heternally1',
// }

// var name = 'zl'

export default 1

// foo.call(obj)
// @ts-ignore
foo.call()
// 如果call、apple、bind的绑定对象是null或者undefined，那么实际上在调用时这些值都会被忽略，所以使用的是默认绑定规则
