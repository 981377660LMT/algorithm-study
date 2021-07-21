"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// 找出该数组中满足其和 ≥ target 的长度最小的 连续子数组
// 如果你已经实现 O(n) 时间复杂度的解法,
// 请尝试设计一个 O(n log(n)) 时间复杂度的解法。
const minSubArrayLen = (target, nums) => {
    let res = 0;
    let l = 0;
    let r = 0;
    let sum = 0;
    // find start
    for (let i = 0; i < nums.length; i++) {
        sum += nums[i];
        if (sum >= target) {
            r = i;
            res = i + 1;
            break;
        }
    }
    if (sum < target)
        return 0;
    while (r < nums.length) {
        l++;
        sum -= nums[l - 1];
        while (sum < target) {
            r++;
            sum += nums[r];
        }
        res = Math.min(res, r - l + 1);
    }
    return res;
};
console.log(minSubArrayLen(7, [2, 3, 1, 2, 4, 3]));
