// console.log(foo) // ReferenceError: foo is not defined

if (
  function foo() {
    console.log('BFE')
  }
) {
  console.log('dev')
}
foo()

// dev
// ReferenceError: foo is not defined
