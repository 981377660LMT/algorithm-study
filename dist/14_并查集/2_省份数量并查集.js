"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const _0____1 = require("./0_\u5E76\u67E5\u96C6");
const findCircleNum = (isConnected) => {
    const uf = new _0____1.UnionFind();
    for (let i = 0; i < isConnected.length; i++) {
        uf.add(i);
        for (let j = 0; j < isConnected.length; j++) {
            if (isConnected[i][j] === 1) {
                uf.union(i, j);
            }
        }
    }
    console.log(uf);
    return uf.count;
};
console.log(findCircleNum([
    [1, 1, 0],
    [1, 1, 0],
    [0, 0, 1],
]));
console.log(findCircleNum([
    [1, 0, 0, 1],
    [0, 1, 1, 0],
    [0, 1, 1, 1],
    [1, 0, 1, 1],
]));
