function foo(a, b, undefined, undefined) {
  console.log('BFE.dev') // this line will not exceute since we never called the function foo
}
console.log(foo.length)
