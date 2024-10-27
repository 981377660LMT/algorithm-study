/* eslint-disable implicit-arrow-linebreak */
/* eslint-disable no-promise-executor-return */

// reduce 实现串行，map 实现并行

async function runSerially(tasks: (() => Promise<void>)[]) {
  return tasks.reduce((pre, cur) => pre.then(cur), Promise.resolve())
}

async function runInParallelly(tasks: (() => Promise<void>)[]) {
  return tasks.map(task => task())
}

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  const tasks = [
    () =>
      new Promise<void>(resolve =>
        setTimeout(() => {
          console.log('task1')
          resolve()
        }, 1000)
      ),
    () =>
      new Promise<void>(resolve =>
        setTimeout(() => {
          console.log('task2')
          resolve()
        }, 2000)
      ),
    () =>
      new Promise<void>(resolve =>
        setTimeout(() => {
          console.log('task3')
          resolve()
        }, 3000)
      )
  ]

  runSerially(tasks).then(() => console.log('runSerially done'))
  runInParallelly(tasks).then(() => console.log('runInParallel done'))
}
