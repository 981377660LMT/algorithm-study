"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
class TreeNode {
    constructor(val = 0, left = null, right = null) {
        this.val = val;
        this.left = left;
        this.right = right;
    }
}
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
const bt2 = {
    val: 1,
    left: { val: 2, left: null, right: null },
    right: null,
};
// 层序遍历后再处理即可
const zipzagLevelOrder = (root) => {
    if (!root)
        return [];
    const res = [];
    const queue = [[root, 0]];
    while (queue.length) {
        const [head, height] = queue.shift();
        if (!res[height]) {
            res[height] = [head.val];
        }
        else {
            res[height].push(head.val);
        }
        head.left && queue.push([head.left, height + 1]);
        head.right && queue.push([head.right, height + 1]);
    }
    return res.map((arr, index) => (index % 2 === 1 ? arr.reverse() : arr));
};
console.log(zipzagLevelOrder(bt));
