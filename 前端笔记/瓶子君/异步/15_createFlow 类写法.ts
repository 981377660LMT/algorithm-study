class Flow {
  private queue: (Func | Flow)[]

  constructor(flows: (Func | Flow)[]) {
    this.queue = flows
  }

  async run(cb?: Func) {
    for (const task of this.queue) {
      if (task instanceof Flow) {
        await task.run()
      } else {
        await task()
      }
    }
    cb && cb()
  }
}

type Func = (...args: any[]) => any
type FuncOrFlow = Func | Flow
type FlowProps = FuncOrFlow | Array<FuncOrFlow>
// createFlow 以一个数组作为参数，数组参数可有以下几种类型:普通函数/异步函数/嵌套createFlow/数组
function createFlow(effects: FlowProps[]): Flow {
  return new Flow(effects.flat())
}

// 异步串行
// 需要按照 a,b,延迟1秒,c,延迟1秒,d,e, done 的顺序打印
const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))
const subFlow = createFlow([() => delay(1000).then(() => console.log('c'))])
createFlow([
  () => console.log('a'),
  () => console.log('b'),
  subFlow,
  [() => delay(1000).then(() => console.log('d')), () => console.log('e')],
]).run(() => {
  console.log('done')
})

export {}
