"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const graph = {
    0: [1, 2],
    1: [2],
    2: [0, 3],
    3: [3],
};
const visited = new Set();
// const dfs = (n: number) => {
//   if (visited.has(n)) return
//   visited.add(n)
//   console.log(n)
//   // 多了一步
//   graph[n].forEach(c => !visited.has(c) && dfs(c))
// }
// 注意:bfs的初始化queue有时候会搭配度排序使用
// 见课程表二的bfs
const bfs = (n) => {
    if (visited.has(n))
        return;
    visited.add(n);
    console.log(n, 6666);
    const queue = [n];
    while (queue.length) {
        const head = queue.shift();
        visited.add(head);
        console.log(head);
        graph[head].forEach(c => {
            if (!visited.has(c)) {
                queue.push(c);
            }
        });
    }
};
// Object.keys(graph).forEach(k => !visited.has(Number(k)) && dfs(Number(k)))
// dfs(2)
bfs(0);
