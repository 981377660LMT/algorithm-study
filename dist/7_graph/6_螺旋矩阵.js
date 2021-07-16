"use strict";
// 给你一个 m 行 n 列的矩阵 matrix ，请按照 顺时针螺旋顺序 ，返回矩阵中的所有元素。
Object.defineProperty(exports, "__esModule", { value: true });
const spiralOrder = (matrix) => {
    const res = [];
    const m = matrix.length;
    const n = matrix[0].length;
    // 记录每一个点走过与否;0/1
    const routesMap = Array.from({ length: m }, () => Array(n).fill(0));
    // 记录是否结束
    let count = 0;
    const getDeltaFromState = (state) => {
        switch (state) {
            case 0:
                return [0, 1];
            case 1:
                return [1, 0];
            case 2:
                return [0, -1];
            case 3:
                return [-1, 0];
            default:
                return [0, 0];
        }
    };
    const dfs = (r, c, state) => {
        if (count === m * n)
            return;
        count++;
        // 走过了 切换
        routesMap[r][c] = 1;
        console.log(r, c);
        res.push(matrix[r][c]);
        const [dR, dC] = getDeltaFromState(state);
        const nextR = r + dR;
        const nextC = c + dC;
        const isNextInMatrix = nextR >= 0 && nextR < m && nextC >= 0 && nextC < n;
        const isNextVisited = isNextInMatrix && routesMap[nextR][nextC] === 1;
        if (isNextInMatrix && !isNextVisited) {
            dfs(nextR, nextC, state);
        }
        else if (!isNextInMatrix || isNextVisited) {
            const newState = ++state % 4;
            const [dR, dC] = getDeltaFromState(newState);
            dfs(r + dR, c + dC, newState);
        }
    };
    dfs(0, 0, 0);
    return res;
};
console.log(spiralOrder([
    [1, 2, 3, 4],
    [5, 6, 7, 8],
    [9, 10, 11, 12],
]));
