"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
let arr = [
    { id: 1, name: '部门1', pid: 0, children: [] },
    { id: 2, name: '部门2', pid: 0, children: [] },
    { id: 3, name: '部门3', pid: 0, children: [] },
    { id: 4, name: '部门4', pid: 1, children: [] },
    { id: 5, name: '部门5', pid: 3, children: [] },
    { id: 6, name: '部门6', pid: 2, children: [] },
    { id: 7, name: '部门7', pid: 5, children: [] },
    { id: 8, name: '部门8', pid: 4, children: [] },
    { id: 9, name: '部门9', pid: 1, children: [] },
    { id: 10, name: '部门10', pid: 2, children: [] },
];
const pidToNodeId = new Map();
// 在graph中建立节点与子节点的关系
arr.forEach(node => pidToNodeId.set(node.pid, [...(pidToNodeId.get(node.pid) || []), node.id]));
const dfs = (nodeId, p) => {
    pidToNodeId.has(nodeId) &&
        pidToNodeId.get(nodeId).forEach(childNodeId => {
            // 根据节点id找到对应的节点
            const childNode = arr[childNodeId - 1];
            p.children.push(childNode);
            dfs(childNodeId, childNode);
        });
};
// 各个根节点开始dfs
const roots = arr.filter(node => node.pid === 0);
roots.forEach(root => dfs(root.id, root));
console.log(pidToNodeId);
console.dir(roots, { depth: null });
