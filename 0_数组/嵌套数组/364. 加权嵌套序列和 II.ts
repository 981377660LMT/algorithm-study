/**
 * // This is the interface that allows for creating nested lists.
 * // You should not implement it, or speculate about its implementation
 * class NestedInteger {
 *     If value is provided, then it holds a single integer
 *     Otherwise it holds an empty nested list
 *     constructor(value?: number) {
 *         ...
 *     };
 *
 *     Return true if this NestedInteger holds a single integer, rather than a nested list.
 *     isInteger(): boolean {
 *         ...
 *     };
 *
 *     Return the single integer that this NestedInteger holds, if it holds a single integer
 *     Return null if this NestedInteger holds a nested list
 *     getInteger(): number | null {
 *         ...
 *     };
 *
 *     Set this NestedInteger to hold a single integer equal to value.
 *     setInteger(value: number) {
 *         ...
 *     };
 *
 *     Set this NestedInteger to hold a nested list and adds a nested integer elem to it.
 *     add(elem: NestedInteger) {
 *         ...
 *     };
 *
 *     Return the nested list that this NestedInteger holds,
 *     or an empty list if this NestedInteger holds a single integer
 *     getList(): NestedInteger[] {
 *         ...
 *     };
 * };
 */
//  令 maxDepth 是任意整数的 最大深度 。
//  整数的 权重 为 maxDepth - (整数的深度) + 1 => 出现得越早，加的越多
//  将 nestedList 列表中每个整数先乘权重再求和，返回该加权和。

// 不用求得深度，只要一层一层的累加进res即可:出现得越早，加的越多；最里层只加最后一次
function depthSumInverse(nestedList: NestedInteger[]): number {
  let res = 0
  let levelSum = 0

  while (nestedList.length > 0) {
    const nextQueue: NestedInteger[] = []

    for (const listOrInteger of nestedList) {
      if (listOrInteger.isInteger()) {
        levelSum += listOrInteger.getInteger()!
      } else {
        nextQueue.push(...listOrInteger.getList()!)
      }
    }

    nestedList = nextQueue
    res += levelSum
  }

  return res
}

export {}

// 输入：nestedList = [[1,1],2,[1,1]]
// 输出：8
// 解释：4 个 1 在深度为 1 的位置， 一个 2 在深度为 2 的位置。
// 1*1 + 1*1 + 2*2 + 1*1 + 1*1 = 8
