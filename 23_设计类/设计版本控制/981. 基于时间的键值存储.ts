class Node {
  constructor(public key: string, public value: string, public timestamp: number) {}
}

/**
 *
 */
class TimeMap {
  private map: Map<string, Node[]>
  constructor() {
    this.map = new Map()
  }

  // set 操作中的时间戳都是严格递增的
  set(key: string, value: string, timestamp: number): void {
    !this.map.has(key) && this.map.set(key, [])
    this.map.get(key)!.push(new Node(key, value, timestamp))
  }

  // 返回对应最大的  timestamp_prev 的那个值   对应bisectRight
  get(key: string, timestamp: number): string {
    const nodeList = this.map.get(key)
    if (!nodeList || !nodeList.length) return ''
    let l = 0
    let r = nodeList.length - 1
    while (l <= r) {
      const mid = (l + r) >> 1
      const node = nodeList[mid]
      if (node.timestamp > timestamp) r = mid - 1
      else if (node.timestamp < timestamp) l = mid + 1
      else return node.value
    }
    return r >= 0 ? nodeList[r].value : ''
  }
}

export {}

// 没有删除操作：哈希表套数组
// 有删除操作：哈希表套树
// 考虑在原题的基础上，增加一个 String del(String k, int t) 的功能：将严格等于键和时间戳的 KV 对删掉。
// 直接使用基于红黑树实现的 TreeMap 即可
