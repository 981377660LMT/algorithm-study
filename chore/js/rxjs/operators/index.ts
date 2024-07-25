// https://rxjs-cn.github.io/learn-rxjs-operators/

import {
  /** Constants */
  EMPTY,
  NEVER,

  /** Internal */
  Observable,

  /** Static observable creation */
  of,
  range,
  from,
  fromEvent,
  interval,
  timer,
  zip,
  concat,

  /** Operators */
  audit,
  auditTime,
  buffer,
  bufferCount,
  bufferTime,
  bufferToggle,
  catchError,
  combineLatest,
  combineLatestAll,
  bufferWhen,
  map,
  take,
  takeWhile,
  distinctUntilChanged,
  distinctUntilKeyChanged,
  forkJoin,
  mergeMap,
  takeUntil,
  scan,
  switchMap,
  switchAll,
  switchScan,
  delay,
  merge,
  fromEventPattern
} from 'rxjs'
import { concatAll, concatWith, mergeAll, pairwise, withLatestFrom } from 'rxjs/operators'

// xxUntil：直到xx操作发生，结束；
// xxWhen：当xx操作发生时，触发；
// xxWhile：当xx操作发生时，继续。

// 创建一个发出随机值的 Observable
const createRandomEmitObservable = () =>
  new Observable<number>(subscriber => {
    let count = 0
    const emitRandomValue = () => {
      count++
      if (count > 10) {
        subscriber.complete()
        return
      }
      const randomValue = Math.floor(Math.random() * 10) // 发出一个 0 到 9 之间的随机整数
      subscriber.next(randomValue)
      const randomDelay = Math.random() * 1000 + 100 // 设置下一次发出值的随机时间间隔（1000ms 到 5000ms之间）
      if (!subscriber.closed) {
        setTimeout(emitRandomValue, randomDelay) // 如果订阅者还没有关闭（即没有取消订阅），则继续递归调用
      }
    }
    emitRandomValue()
    return () => {
      console.log('Observable 已取消订阅')
    }
  })

// ---------------------------- 组合 ----------------------------
// combine、concat、merge
// forkJoin、race
// withLatestFrom、zip、startWith、pairwise

// combineLatest
// 将多个 Observable 合并成一个新的 Observable。
// !这个新的 Observable 会发出一个`数组`，数组中包含了每个输入 Observable 的最新值。
// 只有当所有输入 Observable 中至少每个都发出了至少一个值之后，combineLatest 才会开始发出值
function testCombineLatest() {
  const source = createRandomEmitObservable().pipe(take(2))
  const source2 = createRandomEmitObservable().pipe(take(2))
  const combined = combineLatest([source, source2])
  combined.subscribe({
    next: value => console.log(value),
    complete: () => console.log('Completed')
  })
}

function testCombineLatestAll() {
  // 每秒发出值，并只取前2个
  const source = interval(1000).pipe(take(2))
  // 将 source 发出的每个值映射成取前5个值的 interval observable
  const example = source.pipe(
    map(v =>
      interval(1000).pipe(
        map(i => `Result (${v}): ${i}`),
        take(5)
      )
    )
  )

  /*
  soure 中的2个值会被映射成2个(内部的) interval observables，
  这2个内部 observables 每秒使用 combineLatest 策略来 combineAll，
  每当任意一个内部 observable 发出值，就会发出每个内部 observable 的最新值。
*/
  const combined = example.pipe(combineLatestAll())
  combined.subscribe(console.log)
}

// 你可以把 concat 想象成 ATM 机前的长队，下一次交易 (subscription) 不能在前一个交易完成前开始
// !串行
function testConcat() {
  // 发出 1,2,3
  const sourceOne = of(1, 2, 3)
  // 发出 4,5,6
  const sourceTwo = of(4, 5, 6)

  // 延迟3秒，然后发出
  const sourceThree = sourceOne.pipe(delay(3000))
  // sourceTwo 要等待 sourceOne 完成才能订阅
  const example = sourceThree.pipe(concatWith(sourceTwo))
  // 输出: 1,2,3,4,5,6
  const subscribe = example.subscribe(val => console.log('Example: Delayed source one:', val))

  of(1, 2, 3)
    .pipe(map(val => of(val * 10).pipe(delay(500))))
    .pipe(concatAll())
    .subscribe(console.log)
}

// 将多个 observables 转换成单个 observable
// !并行
function testMerge() {
  const source1 = interval(1000)
  const source2 = interval(2000)
  const res = merge(source1, source2)
  res.subscribe(console.log)

  of(1, 2, 3)
    .pipe(map(val => of(val * 10).pipe(delay(500))))
    .pipe(mergeAll())
    .subscribe(console.log)
}

// 还提供另一个 observable 的最新值
function testWithLatestFrom() {
  const source = interval(5000)
  const secondSource = interval(1000)
  const example = source.pipe(
    withLatestFrom(secondSource),
    map(([first, second]) => {
      return `First Source (5s): ${first} Second Source (1s): ${second}`
    })
  )
  example.subscribe(console.log)
}

function testZip() {
  const source = interval(1000)
  const secondSource = interval(2000)
  const example = zip(source, secondSource)
  example.subscribe(console.log)
}

// startsWith：将一个值或多个值插入到源 Observable 的开头

// ---------------------------- 条件 ----------------------------
// defaultIfEmpty：如果源 Observable 没有发出值，就发出一个默认值
// every：判断所有的值是否都满足条件
// ---------------------------- 创建 ----------------------------
// new、empty、from、fromEvent、fromEventPattern
// interval、of、range、throwError、timer

function testCreate() {
  // 创建一个发出 1,2,3 的 Observable
  const source = new Observable<number>(subscriber => {
    let value = 0
    const interval = setInterval(() => {
      if (value % 2 === 0) {
        subscriber.next(value)
      }
      value++
    }, 1000)

    return () => clearInterval(interval)
  })

  // 输出: 0...2...4...6...8
  const subscribe = source.subscribe(val => console.log(val))
  // 10秒后取消订阅
  setTimeout(() => {
    subscribe.unsubscribe()
  }, 10000)
}

function testTimer() {
  timer(1000, 1000).subscribe(console.log) // startDue、intervalDuration
}

// https://rxjs-cn.github.io/learn-rxjs-operators/operators/creation/frompromise.html
// TODO
function testFromEventPattern() {
  // const addHandler = (handler: Function) => {
  //   document.addEventListener('click', handler)
  // }
  // const removeHandler = (handler: Function) => {
  //   document.removeEventListener('click', handler)
  // }
  // const source = fromEventPattern(addHandler, removeHandler)
  // source.subscribe(console.log)
}

// ---------------------------- 创建 ----------------------------

// audit/auditTime （audit: 审计，auditTime: 审计时间）
// 根据提供的持续时间选择器函数来延迟源Observable（被观察对象）发出的值的发射
// 对于源Observable发出的每个值，audit操作符会使用durationSelector函数来生成一个新的Observable
function testAudit() {
  const source = interval(0)
  const result = source.pipe(audit(v => (v % 2 === 0 ? interval(1000) : of(0))))
  // const result = source.pipe(auditTime(1000))
  result.subscribe(x => console.log(x))
}

// buffer 操作符收集源 Observable 的值，直到传入的闭包 Observable 发出值，然后它发出这些值作为数组
// bufferCount 操作符收集源 Observable 的值，直到收集到一定数量的值，然后发出这些值作为数组
// bufferTime 操作符收集源 Observable 的值，直到过了一定时间，然后发出这些值作为数组。
// bufferToggle 操作符从源 Observable 中收集值，开始收集是由 openings Observable 决定的，结束收集是由 closingSelector 函数决定的
// bufferWhen 操作符收集源 Observable 的值，直到 closingSelector 函数返回的 Observable 发出值，然后发出这些值作为数组
function testBuffer() {
  const source = interval(1000)
  // const buffered = source.pipe(bufferCount(3), bufferTime(100))
  const buffered = source.pipe(bufferWhen(() => interval(1000 + Math.random() * 4000)))
  buffered.subscribe(x => console.log(x))
}

/**
 * !这里面的实现很牛.
 */
function testCatchError() {
  const source = new Observable(subscriber => {
    throw new Error('Oops!')
  })
  source.pipe(catchError(err => EMPTY)).subscribe(console.log)
}

function testDistinctUntilChanged() {
  const source = range(1, 10)
  const result = source.pipe(distinctUntilChanged())
  result.subscribe(console.log)
}

function testTakeUntil() {
  const source = interval(1000)
  const notifier = timer(5000)
  const result = source.pipe(takeUntil(notifier))
  result.subscribe(console.log)
}

// scan 会发出每次累加后的中间值
function testScan() {
  const source = interval(1000).pipe(
    take(5),
    scan((acc, value) => acc + value, 0)
  )
  source.subscribe(console.log)
}

/**
 * !SwitchMap 的特点是，每当源 Observable 发出一个新值时，它会取消订阅之前的内部 Observable 并订阅新的内部 Observable。
 * switchMap 是处理诸如搜索自动完成等需要快速响应最新数据的场景的理想选择，因为它确保了总是只获取最新请求的结果.
 * !switch 意味着`切换到最新`.
 */
function testSwitchMap() {
  // 立即发出值， 然后每5秒发出值
  const source = timer(0, 5000)
  // 当 source 发出值时切换到新的内部 observable，发出新的内部 observable 所发出的值
  const example = source.pipe(switchMap(() => interval(500)))
  // 输出: 0,1,2,3,4,5,6,7,8,9...0,1,2,3,4,5,6,7,8
  const subscribe = example.subscribe(val => console.log(val))
}

/**
 * switchAll 会订阅最新的内部 Observable 并发出它的值，而不是订阅内部 Observable 的 Observable.
 */
function testSwitchAll() {
  const source = interval(1000).pipe(take(3))
  const result = source.pipe(
    map(v => interval(1000).pipe(take(2))),
    switchAll()
  )
  result.subscribe(console.log)
}

function testSwitchScan() {}

/**
 * 类似于 Promise.all，forkJoin 会等待所有的 Observable 都发出值后再发出值
 */
function testForkJoin() {
  const observable1 = of(1).pipe(delay(1000))
  const observable2 = of(2).pipe(delay(2000))
  forkJoin([observable1, observable2]).subscribe({
    next: value => console.log(value),
    complete: () => console.log('Completed')
  })
}

// 组合
// testCombineLatest()
// testConcat()
// testMerge()
// testWithLatestFrom()
// testZip()

// 条件

// 创建
// testCreate()
// testTimer()
testFromEventPattern()

// testAudit()
// testBuffer()
// testCatchError()

// testCombineLatestAll()
// testDistinctUntilChanged()
// testTakeUntil()
// testScan()
// testSwitchMap()
// testSwitchAll()
// testSwitchScan()
// testForkJoin()
