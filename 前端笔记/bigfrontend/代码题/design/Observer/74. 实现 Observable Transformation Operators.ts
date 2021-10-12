import { Observable } from './57. 实现Observable'
import { from } from './70. 实现Observable.from()'

// Observable里有很多的operators，如果我们把Observable想象为一个event stream，
// 那么修改这个stream就成为了常见的需求，transformation operators就是为此而生。
/**
 * @param {any} input
 * @return {(observable: Observable) => Observable}
 * 它把value map到另外一个value，从而生成一个新的event stream。
 */
function map(transform: Function): (observable: Observable) => Observable {
  // your code here
  return source =>
    new Observable(subscriber => {
      // 新的流的subscribe要调用source的subscribe
      source.subscribe(val => subscriber.next(transform(val)))
    })
}

const source = from([1, 2, 3])
map((x: number) => x * x)(source).subscribe(console.log)

// Observable有一个pipe() 方法可以让其更具可读性。

// const source = Observable.from([1,2,3])
// source.pipe(map(x => x * x))
//  .subscribe(console.log)
