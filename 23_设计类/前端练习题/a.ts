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
  private readonly _children: Record<PropertyKey, [handler: HandlerNode, proxy: Obj]> =
    Object.create(null)
  private readonly _action: Record<PropertyKey, Obj> = Object.create(null)

  constructor(obj: Obj) {
    this._origin = obj
  }

  has(_: Obj, prop: PropertyKey): boolean {
    return prop in this._origin || prop in this._action
  }

  get(_: Obj, prop: PropertyKey): Obj {
    console.log(prop, 999)
    if (prop in this._action) return this._action[prop]
    if (prop in this._children) return this._children[prop][1]
    const res = this._origin[prop]
    if (!isObj(res)) return res
    const handler = new HandlerNode(res)
    const proxy = new Proxy(res, handler)
    this._children[prop] = [handler, proxy]
    return proxy
  }

  set(_: Obj, prop: PropertyKey, value: Obj): boolean {
    console.log(prop, 888)
    delete this._action[prop]
    this._action[prop] = value
    return true
  }

  query(): [dirty: boolean, res: Obj] {
    const patch: Record<PropertyKey, Obj> = Object.create(null)
    let dirty = this.dirty
    const keys = Object.keys(this._children)
    keys.forEach(key => {
      const [childDirty, childRes] = this._children[key][0].query()
      console.log(patch, childDirty, childRes, 999)
      if (childDirty) {
        dirty = true
        patch[key] = childRes
      }
    })

    if (!dirty) return [false, this._origin]
    return [true, { ...this._origin, ...patch, ...this._action }]
  }

  get dirty(): boolean {
    return !!Object.keys(this._action).length
  }
}

function isObj(obj: unknown): obj is Obj {
  return typeof obj === 'object' && obj !== null
}
if (require.main === module) {
  const immer = new ImmutableHelper({ arr: [1, 2, 3] })
  const res = immer.produce(draft => (draft.arr[0] = 4))
  console.log(res.arr[0])
}
