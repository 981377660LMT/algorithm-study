"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// 使用Uint8Array加速
const maxVowels = (s, k) => {
    // const weight = new Uint8Array(123)
    // let max = 0
    // weight[97] = weight[101] = weight[105] = weight[111] = weight[117] = 1
    // for (let i = 0; i < k; ++i) {
    //   max += weight[s.charCodeAt(i)]
    // }
    // for (let i = 0, cur = max; i < s.length - k; ++i) {
    //   cur += weight[s.charCodeAt(i + k)] - weight[s.charCodeAt(i)]
    //   cur > max && (max = cur)
    // }
    // return max
    // 使用位运算减少空间复杂度
    const BASE = 97;
    const weight = 1065233; // 1 + (1 << 4) + (1 << 8) + (1 << 14) + (1 << 20);
    let max = 0;
    for (let i = 0; i < k; ++i) {
        max += (weight >> (s.charCodeAt(i) - BASE)) & 1;
    }
    for (let i = 0, cur = max; i < s.length - k; ++i) {
        cur +=
            ((weight >> (s.charCodeAt(i + k) - BASE)) & 1) - ((weight >> (s.charCodeAt(i) - BASE)) & 1);
        cur > max && (max = cur);
    }
    return max;
};
console.log(maxVowels('leetcode', 3));
