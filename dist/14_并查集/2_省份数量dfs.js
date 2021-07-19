"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// 无向图连通域数量
// BFS / DFS 的时间复杂度是 O(n^2)
// 思路:将邻接矩阵初始化为邻接表，对每个点dfs,用visited集合记录走过的点
// 如果visited的size发生了变化，说明产生了新的分支
const findCircleNum = (isConnected) => {
    let res = 0;
    /**
     * @description 记录visited的大小变化
     */
    let preSize = 0;
    const visited = new Set();
    // 初始化邻接表
    const graph = new Map();
    isConnected.forEach((r, ri) => graph.set(ri + 1, r.map((c, ci) => (c === 1 && ri !== ci ? ci + 1 : 0)).filter(v => v !== 0)));
    const dfs = (n) => {
        if (visited.has(n))
            return;
        visited.add(n);
        const v = graph.get(n);
        v.forEach(newN => !visited.has(newN) && dfs(newN));
    };
    for (const n of graph.keys()) {
        // 加不加visited.has对结果影响很大
        !visited.has(n) && dfs(n);
        if (visited.size > preSize) {
            res++;
            preSize = visited.size;
        }
    }
    console.table(graph);
    return res;
};
console.log(findCircleNum([
    [1, 1, 0],
    [1, 1, 0],
    [0, 0, 1],
]));
console.log(findCircleNum([
    [1, 0, 0, 1],
    [0, 1, 1, 0],
    [0, 1, 1, 1],
    [1, 0, 1, 1],
]));
