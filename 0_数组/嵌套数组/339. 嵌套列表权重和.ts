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
//  嵌套列表 [1,[2,2],[[3],2],1] 中每个整数的值就是其深度。
// 请返回该列表按深度加权后所有整数的总和。
function depthSum(nestedList: NestedInteger[]): number {
  return dfs(nestedList, 1)

  // 想象一下二叉树的dfs
  function dfs(nestedList: NestedInteger[], depth: number): number {
    let res = 0

    for (const listOrInteger of nestedList) {
      if (listOrInteger.isInteger()) {
        res += listOrInteger.getInteger()! * depth
      } else {
        res += dfs(listOrInteger.getList(), depth + 1)
      }
    }

    return res
  }
}

export {}

// console.log(depthSum([[1, 1], 2, [1, 1]] ))
// 输出：10
// 解释：因为列表中有四个深度为 2 的 1 ，和一个深度为 1 的 2。
