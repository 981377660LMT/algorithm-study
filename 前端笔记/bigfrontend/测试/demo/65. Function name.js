var foo = function bar() {
  return 'BFE'
}

console.log(foo, foo.name) // [Function: bar] bar
console.log(foo()) // BFE
console.log(bar()) // bar is not defined
