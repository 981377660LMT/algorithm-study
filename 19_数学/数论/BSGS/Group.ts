export {}

/**
 * 群（Group）是代数结构的一种，它包括一个集合和一个二元运算。
 * 群需要满足以下四个性质：
 * 1. 闭合性：对于群中的任意两个元素，其运算结果仍然是群中的元素。
 * 2. 结合性：对于群中的任意三个元素，无论如何结合，其运算结果仍然是群中的元素。
 * 3. 存在单位元：群中存在一个元素，对于群中的任意元素，与其运算都能得到另一个群中的元素。
 * 4. 存在逆元：群中的每个元素都有一个与之对应的元素，对于群中的任意元素，与其运算都能得到单位元。
 */
interface Group<E> {
  e(): E
  op(e1: E, e2: E): E
  inv(e: E): E
  readonly commutative?: boolean
}

/**
 * 半群（SemiGroup）是代数结构的一种，它包括一个集合和一个二元运算。
 * 半群需要满足以下两个性质：
 * 1. 闭合性：对于半群中的任意两个元素，其运算结果仍然是半群中的元素。
 * 2. 结合性：对于半群中的任意三个元素，无论如何结合，其运算结果仍然是半群中的元素。
 */
interface SemiGroup<E> {
  op(e1: E, e2: E): E
}

/**
 * 幺半群（Monoid）是代数结构的一种，它包括一个集合和一个二元运算。
 * 幺半群需要满足以下三个性质：
 * 1. 闭合性：对于幺半群中的任意两个元素，其运算结果仍然是幺半群中的元素。
 * 2. 结合性：对于幺半群中的任意三个元素，无论如何结合，其运算结果仍然是幺半群中的元素。
 * 3. 存在单位元：幺半群中存在一个元素，对于幺半群中的任意元素，与其运算都能得到另一个幺半群中的元素。
 */
interface Monoid<E> extends SemiGroup<E> {
  e(): E
}

/**
 * 阿贝尔群是满足交换律的群。
 */
interface AbleGroup<E> extends Group<E> {
  readonly commutative: true
}

// 树状数组维护的代数结构通常是一个可逆的交换半群（Commutative Semigroup）。
//   如果只需要查询前缀, 那么可以使用一个幺半群（Monoid）。
// 线段树维护的代数结构通常是一个幺半群（Monoid）。
// ST表维护的代数结构通常是一个半群（SemiGroup）。

// Promise是自函子上的幺半群

class PromiseMonoid<E> implements Monoid<Promise<E>> {
  e(): Promise<any> {
    return Promise.resolve()
  }

  op(e1: Promise<E>, e2: Promise<E>): Promise<E> {
    return e1.then(() => e2)
  }
}
