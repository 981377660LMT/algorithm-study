"use strict";
/**
 * 判断有向图中不存在环
 */
Object.defineProperty(exports, "__esModule", { value: true });
// dfs: 先构建出图和isVisted集合,再dfs
const canFinish = (numCourses, prerequisites) => {
    if (prerequisites.length === 0)
        return true;
    const graph = new Map();
    // 访问过的检查为无环的节点,加快计算
    const visited = new Set();
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
    const dfs = (course, deps) => {
        const neighbors = graph.get(course);
        if (neighbors) {
            for (const neighbor of neighbors) {
                // 如果曾经访问过这个节点时不存在环，那么从这节点走下去也不会存在环，直接返回false就可以了
                if (visited.has(neighbor))
                    continue;
                // 本次已经看过的依赖,说明依赖有环
                if (deps.has(neighbor))
                    return (canFinish = false);
                dfs(neighbor, deps.add(neighbor));
            }
        }
        // 验证为无环的节点
        visited.add(course);
    };
    for (const course of graph.keys()) {
        if (!canFinish)
            return false;
        dfs(course, new Set());
    }
    return true;
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
console.log(canFinish(4, [
    [1, 0],
    [2, 0],
    [3, 1],
    [3, 2],
]));
