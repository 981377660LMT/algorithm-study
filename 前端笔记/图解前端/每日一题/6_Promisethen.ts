// @ts-ignore
Promise.resolve().then(1).then(console.log(1))
// then方法提供一个自定义的回调函数，
// !若传入非函数，则会`忽略`当前then方法。
