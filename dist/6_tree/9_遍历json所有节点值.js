"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// 深度优先
const json = {
    a: { b: { c: 1 } },
    d: [1, 2],
};
const dfs = (n, memo = []) => {
    console.log(n, memo);
    Object.keys(n).forEach(k => dfs(n[k], memo.concat(k)));
};
dfs(json);
