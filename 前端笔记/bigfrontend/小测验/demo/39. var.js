function foo() {
  console.log(i)
  for (var i = 0; i < 3; i++) {
    console.log(i)
  }
}

foo()
// undefined
// 0
// 1
// 2
