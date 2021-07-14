"use strict";
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
// 快速比较两个json
const isSameTree = (t1, t2) => {
    // 递归的终点
    if (!t1 && !t2)
        return true;
    if (t1 &&
        t2 &&
        t1.val === t2.val &&
        isSameTree(t1.left, t2.left) &&
        isSameTree(t1.right, t2.right)) {
        return true;
    }
    return false;
};
console.log(isSameTree(bt, bt));
