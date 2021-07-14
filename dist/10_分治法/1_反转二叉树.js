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
// 时间复杂 O(树的节点数)
// 空间复杂 O(树的高度)
const reverseBinaryTree = (bt) => {
    const reverseTwo = (bt) => ([bt.left, bt.right] = [bt.right, bt.left]);
    reverseTwo(bt);
    bt.left && reverseBinaryTree(bt.left);
    bt.right && reverseBinaryTree(bt.right);
};
reverseBinaryTree(bt);
console.log(bt);
