// https://leetcode.com/problems/memoize-ii/
// 1 <= inputs.length <= 105
// 0 <= inputs.flat().length <= 105
// inputs[i][j] != NaN

// 基于参数树状结构缓存调用结果，参考自React
// https://leetcode.cn/problems/memoize-ii/solution/ji-yu-can-shu-shu-zhuang-jie-gou-huan-cu-m5l9/

type Fn = (...params: any) => any

const enum CallState {
  // 参数未结束
  Default,
  // 返回值
  Return,
  // 调用出错
  Error
}

interface MemoNode<T extends (...args: unknown[]) => unknown> {
  state: CallState
  value: ReturnType<T> | null
  error: unknown
  // 存储基本类型参数变量
  primitive: Map<unknown, MemoNode<T>>
  // 存储引用类型参数变量（使用WeakMap防止内存泄漏）
  reference: WeakMap<any, MemoNode<T>>
}

function newMemoNode<T extends (...args: unknown[]) => unknown>(): MemoNode<T> {
  return {
    state: CallState.Default,
    value: null,
    error: null,
    primitive: new Map(),
    reference: new WeakMap()
  }
}

function isPrimitiveType(value: unknown) {
  return (typeof value !== 'object' && typeof value !== 'function') || value === null
}

function memoize(fn: Fn): Fn {
  const root = newMemoNode()
  return function (...args: unknown[]) {
    // 判断是否有缓存，有则直接返回，无则新建缓存节点
    let cur = root
    for (const arg of args) {
      if (isPrimitiveType(arg)) {
        if (cur.primitive.has(arg)) {
          cur = cur.primitive.get(arg)!
        } else {
          const newNode = newMemoNode()
          cur.primitive.set(arg, newNode)
          cur = newNode
        }
      } else if (cur.reference.has(arg)) {
        cur = cur.reference.get(arg)!
      } else {
        const newNode = newMemoNode()
        cur.reference.set(arg, newNode)
        cur = newNode
      }
    }

    // 判断当前节点是否缓存结果
    if (cur.state === CallState.Return) {
      return cur.value
    }

    if (cur.state === CallState.Error) {
      throw cur.error
    }

    // 无缓存则执行并创建缓存
    try {
      const value = fn(...args)
      cur.state = CallState.Return
      cur.value = value
      return value
    } catch (error) {
      cur.state = CallState.Error
      cur.error = error
      throw error
    }
  }
}

export {}
