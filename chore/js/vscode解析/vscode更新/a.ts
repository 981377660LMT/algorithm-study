function foo() {
  let result: number
  if (someCondition()) {
    result = doSomeWork()
  } else {
    let temporaryWork = doSomeWork()
    temporaryWork *= 2
    // forgot to assign to 'result'
  }

  printResult()

  function printResult() {
    console.log(result) // no error here.
  }
}

export {}
