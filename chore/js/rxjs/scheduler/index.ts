// 调度器，用于控制计算何时执行
// RxJS 提供了几种内置的 Scheduler 类型，每种类型都有其特定的用途。
// Scheduler 的类型
// null: 如果不指定 Scheduler，操作将默认同步执行。
// asyncScheduler: 使得任务异步执行，类似于使用 setTimeout。
// asapScheduler: 尽快执行任务，但是在当前事件循环结束后，类似于 Promise.resolve()。
// queueScheduler: 在当前事件循环中执行任务，按照队列的方式依次执行。
// animationFrameScheduler: 用于在浏览器的下一个重绘之前执行任务，适合进行动画相关的操作

import { of, asyncScheduler, interval, animationFrameScheduler, SchedulerAction, Scheduler, SchedulerLike } from 'rxjs'
import { observeOn, takeWhile } from 'rxjs/operators'

// 创建一个 Observable，并指定使用 asyncScheduler
// !observeOn 操作符用于指定 observable 在 asyncScheduler 上执行
function testAsyncScheduler() {
  const observable = of(1, 2, 3).pipe(observeOn(asyncScheduler))

  console.log('Before subscribe')

  observable.subscribe({
    next: value => console.log(value),
    complete: () => console.log('Completed')
  })

  console.log('After subscribe')
}

function testAnimationFrameScheduler() {
  // 使用 animationFrameScheduler 创建一个 Observable
  // !它在每个浏览器动画帧上发出值，直到值达到 60
  const animationObservable = interval(0, animationFrameScheduler).pipe(takeWhile(val => val < 60))

  animationObservable.subscribe({
    next: val => console.log(`Animation frame: ${val}`),
    complete: () => console.log('Animation completed')
  })
}

if (require.main === module) {
  // testAsyncScheduler()
  // testAnimationFrameScheduler()
}
