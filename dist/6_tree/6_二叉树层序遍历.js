"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const bt = {
    val: 1,
    left: {
        val: 2,
        left: {
            val: 4,
            left: null,
            right: null,
        },
        right: {
            val: 5,
            left: null,
            right: null,
        },
    },
    right: {
        val: 3,
        left: {
            val: 6,
            left: null,
            right: null,
        },
        right: {
            val: 7,
            left: null,
            right: null,
        },
    },
};
// 返回值的第h层是一个装满该层节点值的数组number[][]
const getLevelOrder = (root) => {
    if (!root)
        return [];
    const queue = [[root, 0]];
    const res = [];
    while (queue.length) {
        const [head, level] = queue.shift();
        if (!res[level]) {
            res[level] = [head.val];
        }
        else {
            res[level].push(head.val);
        }
        console.log(head.val, level);
        head.left && queue.push([head.left, level + 1]);
        head.right && queue.push([head.right, level + 1]);
    }
    return res;
};
console.log(getLevelOrder(bt));
