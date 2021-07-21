"use strict";
// 多数元素是指在数组中出现次数 大于 ⌊ n/2 ⌋ 的元素。
// 数组是非空的，并且给定的数组总是存在多数元素。
// 设计时间复杂度为 O(n)、空间复杂度为 O(1) 的算法解决此问题。
Object.defineProperty(exports, "__esModule", { value: true });
const majorityElement = (nums) => nums.sort()[Math.floor(nums.length / 2)];
console.log(majorityElement([2, 2, 1, 1, 1, 2, 2]));
