"use strict";
// 其实用map更好O(n)
Object.defineProperty(exports, "__esModule", { value: true });
// 时间复杂度O(n^2)
const intersection = (arr1, arr2) => [...new Set(arr1)].filter(ele => arr2.includes(ele));
console.log(intersection([1, 2, 3], [1]));
