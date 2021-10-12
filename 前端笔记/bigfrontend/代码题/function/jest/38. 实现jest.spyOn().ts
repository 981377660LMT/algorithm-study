// 如果你写过单元测试的话，一定很熟悉Spy的用法。
// 请自己实现一个spyOn(object, methodName) ，类似于 jest.spyOn()。
// 以下是spyOn需要完成的内容。

// 1.spy被调用的时候，原来的method也需要被调用。
// 2.spy需要又一个calls数组，数组中含有所有调用的参数
// Jest中：
// interface Spy {
//   (...params: any[]): any;
//   identity: string;
//   and: SpyAnd;
//   calls: Calls;
//   mostRecentCall: { args: any[] };
//   argsForCall: any[];
//   wasCalled: boolean;
// }

// jasmine 意为茉莉花
interface JasmineSpy {
  calls: any[]
}

/**
 * @param {object} obj
 * @param {string} methodName
 */
function spyOn<T>(object: T, methodName: keyof T): JasmineSpy {
  // your code here
  const method = object[methodName]
  if (typeof method !== 'function') throw new Error(`${methodName} is not function`)

  const calls: any[] = []

  // @ts-ignore
  object[methodName] = function (...args: any[]) {
    calls.push(args)
    method.apply(this, args)
  }

  return { calls }
}

if (require.main === module) {
  const obj = {
    data: 1,
    increment(num: number) {
      this.data += num
    },
  }

  const spy = spyOn(obj, 'increment')
  obj.increment(1)
  console.log(obj.data) // 2
  obj.increment(2)
  console.log(obj.data) // 4
  console.log(spy.calls)
  // [ [1], [2] ]
}

export {}
