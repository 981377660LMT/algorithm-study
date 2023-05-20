/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable prefer-destructuring */

// 创建带有微小修改的不可变对象的克隆副本是一个繁琐的过程。
// 请你编写一个名为 ImmutableHelper 的类，作为满足这一要求的工具。构造函数接受一个不可变对象 obj ，该对象将是一个 JSON 对象或数组。
// 该类有一个名为 produce 的方法，它接受一个名为 mutator 的函数。该函数返回应用这些变化的 obj 的副本。
// mutator 函数接受 obj 的 代理 版本。函数的使用者可以（看起来）对该对象进行修改，但原始对象 obj 实际上没有被改变。

// 例如，用户可以编写如下代码：
// const originalObj = {"x": 5};
// const helper = new ImmutableHelper(originalObj);
// const newObj = helper.produce((proxy) => {
//   proxy.x = proxy.x + 1;
// });
// console.log(originalObj); // {"x": 5}
// console.log(newObj); // {"x": 6}

// !mutator 函数的属性：
//  它始终返回 undefined 。
//  它永远不会访问不存在的键。
//  它永远不会删除键（ delete obj.key ）。
//  它永远不会在代理对象上调用方法（ push 、shift 等）。
//  它永远不会将键设置为对象（ proxy.x = {} ）。
//  关于如何测试解决方案的说明：解决方案验证器仅分析返回结果与原始 obj 之间的差异。进行完全比较的计算开销太大。
// 2 <= JSON.stringify(obj).length <= 4e5
// produce() 的总调用次数 < 1e5

// !1.immer => proxy + copy on write, Cache下diff. 每次修改最坏是O(n).
//    本来复制一遍就这样，如果结合immutable数据结构可以不用O(n).
//    如果修改很深，所有结点全变成dirty, 修改和深拷贝一样了.
//    只适用于很少且很小的改动.
//    大的还是得上数据结构.
// !2.immutable.js => 持久化对象(trie). 每次修改最坏是O(logn).

// https://leetcode.cn/problems/immutability-helper/solution/typescript-proxy-copy-on-write-shi-xian-hsm6b/

type Obj = Record<PropertyKey, unknown> | unknown[]

/**
 * 实现immer库的produce函数.
 */
class ImmutableHelper {
  private readonly _base: Obj

  constructor(draft: Obj) {
    this._base = draft
  }

  produce(recipe: (draft: Obj) => void): Obj {
    const handler = new HandlerNode(this._base)
    recipe(new Proxy(this._base, handler))
    return handler.query()[1]
  }
}

/**
 * const originalObj = {"x": 5};
 * const mutator = new ImmutableHelper(originalObj);
 * const newObj = mutator.produce((proxy) => {
 *   proxy.x = proxy.x + 1;
 * });
 * console.log(originalObj); // {"x: 5"}
 * console.log(newObj); // {"x": 6}
 */

/**
 * 类似于TrieNode的结构.
 * @property _origin 原始对象.
 * @property _children 子节点, 保存每个儿子的handler和proxy.
 * @property _action 每个结点处的修改.
 */
class HandlerNode {
  private readonly _origin: Obj
  private readonly _children: Map<PropertyKey, [handler: HandlerNode, proxy: Obj]> = new Map()
  private readonly _action: Map<PropertyKey, Obj> = new Map()

  constructor(obj: Obj) {
    this._origin = obj
  }

  has(_: Obj, prop: PropertyKey): boolean {
    return prop in this._origin || this._action.has(prop)
  }

  get(_: Obj, prop: PropertyKey): Obj {
    if (this._action.has(prop)) return this._action.get(prop)!
    if (this._children.has(prop)) return this._children.get(prop)![1]
    const res = this._origin[prop]
    if (!isObj(res)) return res
    const handler = new HandlerNode(res)
    const proxy = new Proxy(res, handler)
    this._children.set(prop, [handler, proxy])
    return proxy
  }

  set(_: Obj, prop: PropertyKey, value: Obj): boolean {
    this._children.delete(prop)
    this._action.set(prop, value)
    return true
  }

  query(): [dirty: boolean, res: Obj] {
    const patch: Record<PropertyKey, Obj> = Object.create(null)
    let dirty = this.dirty
    this._children.forEach(([childHandler], key) => {
      const [childDirty, childRes] = childHandler.query()
      if (childDirty) {
        dirty = true
        patch[key] = childRes
      }
    })

    if (!dirty) return [false, this._origin]
    const res = { ...this._origin, ...patch }
    this._action.forEach((value, key) => {
      res[key] = value
    })
    return [true, res]
  }

  get dirty(): boolean {
    return !!this._action.size
  }
}

function isObj(obj: unknown): obj is Obj {
  return typeof obj === 'object' && obj !== null
}

export {}
if (require.main === module) {
  const immer = new ImmutableHelper({ arr: [1, 2, 3] })
  const res = immer.produce(draft => (draft.arr[0] = 4))
  console.log(res.arr[0])
}
