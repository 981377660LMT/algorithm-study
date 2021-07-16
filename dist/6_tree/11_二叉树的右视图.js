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
// 解决方法1:深度遍历并记录高度,先左后右
const rightSideView = (root) => {
    const res = [];
    const dfs = (root, height) => {
        if (!root)
            return;
        res[height] = root.val;
        root.left && dfs(root.left, height + 1);
        root.right && dfs(root.right, height + 1);
    };
    dfs(root, 0);
    return res;
};
// 解决方法2:层序遍历+提取最右边的值
// res的第n层是一个数组
// const rightSideView = (root: TreeNode) => {
//   const res: number[][] = []
//   const queue: [TreeNode, number][] = [[root, 0]]
//   const bfs = (root: TreeNode | null, height: number) => {
//     while (queue.length) {
//       const [head, height] = queue.shift()!
//       if (!res[height]) {
//         res[height] = [head.val]
//       } else {
//         res[height].push(head.val)
//       }
//       console.log(head.val)
//       head.left && queue.push([head.left, height + 1])
//       head.right && queue.push([head.right, height + 1])
//     }
//   }
//   bfs(root, 0)
//   return res.map(arr => arr.slice(-1)).flat()
// }
console.log(rightSideView(bt));
console.log(rightSideView(bt2));
