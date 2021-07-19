"use strict";
// 无向图连通域数量
// BFS / DFS 的时间复杂度是 O(n^2)
Object.defineProperty(exports, "__esModule", { value: true });
// 固定写法
// 注意到isConnected必定是对称矩阵
const findCircleNum = (isConnected) => {
    let circle = 0;
    const visited = new Set();
    const dfs = (i) => {
        for (let j = 0; j < isConnected.length; j++) {
            if (isConnected[i][j] === 1 && !visited.has(j)) {
                visited.add(j);
                dfs(j);
            }
        }
    };
    for (let i = 0; i < isConnected.length; i++) {
        if (!visited.has(i)) {
            dfs(i);
            circle++;
        }
    }
    return circle;
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
