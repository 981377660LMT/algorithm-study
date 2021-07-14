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
// 根左右
const preOrder = (root) => {
    if (!root)
        return [];
    console.log(root.val);
    root.left && preOrder(root.left);
    root.right && preOrder(root.right);
};
// 左根右
const inOrder = (root) => {
    if (!root)
        return [];
    root.left && inOrder(root.left);
    console.log(root.val);
    root.right && inOrder(root.right);
};
// 根右左
const postOrder = (root) => {
    if (!root)
        return [];
    console.log(root.val);
    root.right && postOrder(root.right);
    root.left && postOrder(root.left);
};
// preOrder(bt)
inOrder(bt);
