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
// 都是O(n)
const inorderTraversal = (root) => {
    if (!root)
        return [];
    const res = [];
    const stack = [];
    while (stack.length > 0 || root) {
        if (root) {
            stack.push(root);
            root = root.left;
        }
        else {
            root = stack.pop();
            res.push(root.val);
            root = root.right;
        }
    }
    return res;
};
console.log(inorderTraversal(bt));
