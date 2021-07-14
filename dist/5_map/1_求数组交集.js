"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// 时间复杂度O(n+m)
// 空间复杂度O(m)
const intersection = (arr1, arr2) => {
    const map = new Map();
    const res = [];
    arr1.forEach(n => {
        map.set(n, true);
    });
    arr2.forEach(n => {
        if (map.get(n)) {
            res.push(n);
            map.delete(n);
        }
    });
    return res;
};
console.log(intersection([1, 2, 3], [1]));
