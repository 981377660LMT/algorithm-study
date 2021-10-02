// group[i] 表示第 i 个项目所属的小组，如果这个项目目前无人接手，那么 group[i] 就等于 -1
// 同一小组的项目，排序后在列表中彼此相邻。 -1 表示不属于任一个组
// 项目之间存在一定的依赖关系，我们用一个列表 beforeItems 来表示，其中 beforeItems[i] 表示在进行第 i 个项目前（位于第 i 个项目左侧）应该完成的所有项目。
// n 个项目和  m 个小组
function sortItems(n: number, m: number, group: number[], beforeItems: number[][]): number[] {}

console.log(sortItems(8, 2, [-1, -1, 1, 0, 0, 1, 0, -1], [[], [6], [5], [6], [3, 6], [], [], []]))
// 输出：[6,3,4,1,5,2,0,7]

// 太难了 不做
