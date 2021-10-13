var foo = 1
;(function () {
  // 'var foo = 3' is hoisted here -> var foo = undefined
  console.log(foo) // undefined - local variable
  foo = 2
  console.log(window.foo) // 1 - window's variable
  console.log(foo) // 2 - local variable
  var foo = 3
  console.log(foo) // 3 - local variable
  console.log(window.foo) // 1 - window's variable
})()
