"use strict";
const foo = (matrix) => {
    if (!matrix || !matrix[0])
        return [];
    const m = matrix.length;
    const n = matrix[0].length;
    // m行n列矩阵
    const flow1 = Array.from({ length: m }, () => Array(n).fill(false));
    const flow2 = Array.from({ length: m }, () => Array(n).fill(false));
    const dfs = (r, c, flow) => {
        flow[r][c] = true;
        [
            [r - 1, c],
            [r + 1, c],
            [r, c - 1],
            [r, c + 1],
        ].forEach(([nextR, nextC]) => {
            // 1.在矩阵中
            // 2.没有重复
            // 3.逆流而上
            if (nextR >= 0 &&
                nextR <= m &&
                nextC >= 0 &&
                nextC <= n &&
                !flow[nextR][nextC] &&
                matrix[nextR][nextC] >= matrix[r][c]) {
                dfs(nextR, nextC, flow);
            }
        });
    };
    // 沿着海岸线逆流而上
    for (let r = 0; r < m; r++) {
        dfs(r, 0, flow1);
        dfs(r, n - 1, flow2);
    }
    for (let c = 0; c < n; c++) {
        dfs(0, c, flow1);
        dfs(m - 1, c, flow2);
    }
    // 对比能留到两个大洋里的坐标
    const res = [];
    for (let r = 0; r < m; r++) {
        for (let c = 0; c < n; c++) {
            if (flow1[r][c] && flow2[r][c]) {
                res.push([r, c]);
            }
        }
    }
    return res;
};
// 时间复杂度m*n
// 空间复杂度m*n
