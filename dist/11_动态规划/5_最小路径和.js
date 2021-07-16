"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// 给定一个包含非负整数的 m x n 网格 grid ，请找出一条从左上角到右下角的路径，使得路径上的数字总和为最小。
// O(mn) time, O(1) space
const minPathSum = (grid) => {
    const m = grid.length;
    const n = grid[0].length;
    const minPathSumMatrix = Array.from({ length: m }, () => Array(n).fill(Infinity));
    const firstRow = grid[0];
    const firstColum = grid.map(arr => arr[0]);
    // 初始化
    for (let i = 0; i < n; i++) {
        minPathSumMatrix[0][i] = firstRow.slice(0, i + 1).reduce((pre, cur) => pre + cur, 0);
    }
    for (let i = 0; i < m; i++) {
        minPathSumMatrix[i][0] = firstColum.slice(0, i + 1).reduce((pre, cur) => pre + cur, 0);
    }
    for (let i = 1; i < m; i++) {
        for (let j = 1; j < n; j++) {
            minPathSumMatrix[i][j] =
                Math.min(minPathSumMatrix[i - 1][j], minPathSumMatrix[i][j - 1]) + grid[i][j];
        }
    }
    console.table(minPathSumMatrix);
    return minPathSumMatrix[m - 1][n - 1];
};
console.log(minPathSum([
    [1, 3, 1],
    [1, 5, 1],
    [4, 2, 1],
]));
