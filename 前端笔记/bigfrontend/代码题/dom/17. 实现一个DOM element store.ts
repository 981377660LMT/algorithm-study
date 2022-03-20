// 请在不利用Map的条件下实现一个Node Store，支持DOM element作为key。
// 你可以实现一个通用的Map polyfill。或者利用DOM元素的特性来做文章？

interface StoredNode extends Node {
  $nodeStoreKey: symbol
}

// 思路：使用Symbol来验证身份
// 也可以用node.dataset来 挂载 id
class NodeStore {
  private store: Record<symbol, any> = {}

  /**
   * @param {Node} node
   * @param {any} value
   */
  set(node: StoredNode, value: any) {
    node.$nodeStoreKey = Symbol()
    this.store[node.$nodeStoreKey] = value
  }

  /**
   * @param {StoredNode} node
   * @return {any}
   */
  get(node: StoredNode): any {
    return this.store[node.$nodeStoreKey]
  }

  /**
   * @param {StoredNode} node
   * @return {Boolean}
   */
  has(node: StoredNode): boolean {
    return !!this.store[node.$nodeStoreKey]
  }
}
