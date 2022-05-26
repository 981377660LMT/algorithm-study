import { useUnionFindArray } from '../useUnionFind'

/**
 * @param {number} n
 * @param {number[][]} edgeList  两个点之间可能有 超过一条边 。
 * @param {number[][]} queries
 * @return {boolean[]}
 * 对于每个查询 queries[j] ，判断是否存在从 pj 到 qj 的路径，且这条路径上的每一条边都 严格小于 limitj 。
 * 采取离线排序优化的方式来解。
 * @summary
 * 离线排序
 */
const distanceLimitedPathsExist = function (
  n: number,
  edgeList: number[][],
  queries: number[][]
): boolean[] {
  const res = Array<boolean>(queries.length).fill(false)
  const uf = useUnionFindArray(n)
  edgeList.sort((a, b) => a[2] - b[2])
  queries = queries.map((v, i) => [...v, i]).sort((a, b) => a[2] - b[2])

  let edgeIndex = 0
  queries.forEach(([from, to, limit, index]) => {
    while (edgeIndex < edgeList.length && edgeList[edgeIndex][2] < limit) {
      uf.union(edgeList[edgeIndex][0], edgeList[edgeIndex][1])
      edgeIndex++
    }
    if (uf.isConnected(from, to)) res[index] = true
  })

  return res
}

console.log(
  distanceLimitedPathsExist(
    3,
    [
      [0, 1, 2],
      [1, 2, 4],
      [2, 0, 8],
      [1, 0, 16],
    ],
    [
      [0, 1, 2],
      [0, 2, 5],
    ]
  )
)

// 输出：[false,true]

// 什么叫在线算法？就是依次处理每一个 query，对每一个 query 的计算，
// 和之后的 query 无关，也不会用到之后的 query 信息（但可能也可以使用之前的 query 信息）。

// 所以，在线算法，可以用来处理数据流。
// 算法不需要一次性地把所有的 query 都收集到再处理。
// 大家也可以想象成：把这个算法直接部署到线上，尽管在线上可能又产生了很多新的 query，也不影响，算法照常运行。

// 离线算法则不同。离线算法需要把所有的信息都收集到，才能运行。
// 处理当前 query 的计算过程，可能需要使用之后 query 的信息。

// 插入排序算法是一种在线算法,因为可以把插入排序算法的待排序数组看做是一个数据流。
// 选择排序算法则是一种离线算法。因为选择排序算法一上来要找到整个数组中最小的元素；这就要求不能再有新的数据了。

// 再举一个例子，对于 topK 问题（找前 k 小或者前 k 大的元素）。
// 使用一个大小为 k 的优先队列是在线算法，虽然时间复杂度是 O(nlogk)，但整个算法不需要一次性知道所有的数据，可以处理数据流；
// 使用快排的思想做改进，topK 问题可以在 O(n) 时间解决。但这是一个离线的算法。初始必须知道所有的数据，才能完成。
