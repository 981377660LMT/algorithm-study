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
// 是否为镜像对称的树
const isSymmetric = (root) => {
    const isSymmetricTwo = (r1, r2) => {
        // 递归的终点
        if (!r1 && !r2)
            return true;
        return (!!r1 &&
            !!r2 &&
            r1.val === r2.val &&
            isSymmetricTwo(r1.left, r1.right) &&
            isSymmetricTwo(r1.right, r2.left));
    };
    return isSymmetricTwo(root.left, root.right);
};
console.log(isSymmetric({ val: 1, left: null, right: null }));
