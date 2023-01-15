// 164. 实现Immer的produce()
// !使用 produce 简化深度更新
// https://segmentfault.com/a/1190000042282263
// https://yo-cwj.com/2020/09/16/brief-read-immer/

// !immer本身浅拷贝(写时浅拷贝)还是O(n) 只能方便修改对象(proxy的奇技淫巧) 无法提升性能
// 如果需要提升性能, 可以使用immutable.js
// 如果需要操作方便，可以使用immer.js
// !(很多人误解了可持久化的不可变)
// https://coder-mike.com/blog/2021/03/05/immutable-js-vs-immer/

type ProduceFunc = <T extends Record<PropertyKey, unknown>>(
  base: T,
  recipe: (draft: T) => void
) => object

const produce: ProduceFunc = (base, recipe) => {
  const proxy = createProxy(base) // 生成draft(proxy) 有base和copy两个属性
  recipe(proxy) // mutate draft
  return getResult(proxy) // 分析draft, 取得期望的返回值
}

const NO_PROXY = Symbol('NO_PROXY')
interface State {
  /**
   * 原始数据存在base中，引用方式为地址引用
   */
  base: Record<PropertyKey, unknown>

  /**
   * 修改后的值存在copy中，引用方式为浅拷贝原始数据
   * 如果copy对象更深层次得值被改变,copy对象的值会是一个(递归)proxy对象
   */
  copy: Record<PropertyKey, unknown> | null

  [NO_PROXY]: State
}

function createProxy(origin: Record<PropertyKey, unknown>): State {
  const state: State = {
    base: origin, // 保存原数据
    copy: null
  }

  const handler: ProxyHandler<State> = {
    set(target: State, prop: string, newValue: unknown): boolean {
      target.copy = target.copy ?? { ...target.base } // !浅拷贝(数组很大时O(n)?)
      target.copy[prop] = newValue
      return true
    },
    get(target: State, prop: PropertyKey): State {
      if (prop === NO_PROXY) return target // !STATE_KEY属性不会走拦截
      target.copy = target.copy ?? { ...target.base }
      const nextProxy = createProxy(origin) // 为什么要递归创建proxy?
      target.copy[prop] = nextProxy
      return nextProxy
    }
  }

  // immer 中用 Proxy.revocable 来产生 proxy
  const proxy = new Proxy<State>(state, handler)
  return proxy
}

// 根据draft的特点来获取期望的immutable数据
function getResult(draft: State): State['copy'] {
  const state = draft[NO_PROXY] // 不会走拦截
  for (const key of Object.keys(state.copy!)) {
    const value = state.copy![key]
    if (isDraft(value)) {
      const res = getResult(value)
      state.copy![key] = res
    }
  }
  return state.copy
}

function isDraft(o: unknown): o is State {
  return typeof o === 'object' && o !== null && NO_PROXY in o
}

if (require.main === module) {
  const obj = [
    {
      name: 'BFE'
    },
    {
      name: '.'
    }
  ]
  const newState = produce(obj, draft => {
    draft[0].name = 'bigfrontend'
    draft[0].name = '1'
  })
  console.log(newState[0].name)
}

// 在12. 实现 Immutability helper中, 我们实现了各种helper。
// 但是这些helper需要记住其使用方法，Immer 使用了另外一种方法，可以更加简单。

// 比如我们有如下state。

// const state = [
//   {
//     name: 'BFE',
//   },
//   {
//     name: '.',
//   }
// ]

// 使用produce() 的话，可以按照意愿修改state。

// const newState = produce(state, draft => {
//   draft.push({name: 'dev'})
//   draft[0].name = 'bigfrontend'
//   draft[1].name = '.' // set为相同值。
// })

// 注意，未变化的部分并没有拷贝。
// expect(newState).not.toBe(state);
// expect(newState).toEqual(
//   [
//     {
//       name: 'bigfrontend',
//     },
//     {
//       name: '.',
//     },
//     {
//       name: 'dev'
//     }
//   ]
// );
// expect(newState[0]).not.toBe(state[0])
// expect(newState[1]).toBe(state[1])
// expect(newState[2]).not.toBe(state[2])
// 请实现你的produce().

// 目的并不是重写Immer，test case并不是要cover所有情况。
// !只需要实现简单对象和数组的基本使用方式，Map/Set和Auto freezing等不需要考虑。
// 需要保证未改变的数据部分不被拷贝。

// proxy ??
// 请注意，命名 recipe 的第一个参数 draft 并不是绝对必要的。
// 您可以将其命名为任何您想要的名称，例如 user。
// !使用 draft 作为名称只是一个约定，以表明：“这里的 mutation 是可以的”。

// 原理:代理+懒更新
// 1. Copy on write
// 2. Proxies
// !immer产生immutable对象的原理是递归浅拷贝.
// !利用了proxy来实现懒处理, 没有被touch的对象不会创建新的浅拷贝,
// 依旧使用原对象的内存地址, 节省内存, 提高性能.
