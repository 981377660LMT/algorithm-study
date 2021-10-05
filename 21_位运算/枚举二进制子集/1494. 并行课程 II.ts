// dependencies[i] = [xi, yi]  表示一个先修课的关系，也就是课程 xi 必须在课程 yi 之前上
// 在一个学期中，你 最多 可以同时上 k 门课，前提是这些课的先修课在之前的学期里已经上过了。
// 1 <= n <= 15

/**
 *
 * @param n
 * @param relations
 * @param k
 * 并行算法parallel_algorithms
 */
function minNumberOfSemesters(n: number, relations: number[][], k: number): number {}

console.log(
  minNumberOfSemesters(
    4,
    [
      [2, 1],
      [3, 1],
      [1, 4],
    ],
    2
  )
)
// 在第一个学期中，我们可以上课程 2 和课程 3 。然后第二个学期上课程 1 ，第三个学期上课程 4 。
// n个用时相同的任务（work）有先后依赖关系，
// 现只有k台机器（parallelism 并行度），最少用时多少能完成。最大深度为depth。
