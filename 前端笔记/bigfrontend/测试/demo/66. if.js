if (true) {
  function foo() {
    console.log('BFE')
  }
}
if (false) {
  function bar() {
    console.log('dev')
  }
}

foo() // BFE
bar() // TypeError: bar is not a function
