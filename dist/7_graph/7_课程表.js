"use strict";
// 错误的尝试
// const canFinish = (numCourses: number, prerequisites: [number, number][]) => {
//   if (prerequisites.length === 0) return true
//   // 深度遍历?
//   const map = new Map<number, number>(prerequisites)
//   let canFinish = false
Object.defineProperty(exports, "__esModule", { value: true });
// dfs: 先构建出图和isVisted集合,再dfs
const canFinish = (numCourses, prerequisites) => {
    if (prerequisites.length === 0)
        return true;
    const graph = new Map();
    // 是否访问过
    const visited = new Set();
    // 检查循环
    const visiting = new Set();
    let canFinish = true;
    // 1. 初始化图
    for (const [course, dependentCourse] of prerequisites) {
        if (graph.has(course)) {
            graph.get(course).push(dependentCourse);
        }
        else {
            graph.set(course, [dependentCourse]);
        }
    }
    console.log(graph);
    // 2. dfs遍历图检查是否有没有出口的环
    const dfs = (course) => {
        visiting.add(course);
        const neighbors = graph.get(course);
        if (neighbors) {
            for (const neighbor of neighbors) {
                if (visited.has(neighbor))
                    continue;
                if (visiting.has(neighbor))
                    return (canFinish = false);
                dfs(neighbor);
            }
        }
        visiting.delete(course);
        visited.add(course);
    };
    for (const course of graph.keys())
        dfs(course);
    return canFinish;
};
console.log(canFinish(2, [
    [1, 0],
    [0, 1],
]));
console.log(canFinish(3, [
    [1, 0],
    [0, 1],
    [1, 2],
]));
console.log(canFinish(20, [
    [0, 10],
    [3, 18],
    [5, 5],
    [6, 11],
    [11, 14],
    [13, 1],
    [15, 1],
    [17, 4],
]));
