Promise.resolve(1)
  .then(() => 2) // 2 (as 1 isn't used)
  .then(3) // // skipped, param NOT a function  // Argument of type '3' is not assignable to parameter of type '(value: number) => number | PromiseLike<number>'.ts(2345)
  .then(value => value * 3)
  .then(Promise.resolve(4)) // skipped, param NOT a function // it return an promise object / creates a Pending promise
  .then(console.log) // 6

Promise.resolve(1).then(Promise.resolve(4)).then(console.log) // 1

// skipped, param NOT a function
