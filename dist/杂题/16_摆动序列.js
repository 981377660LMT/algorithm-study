"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// 你能否用 O(n) 时间复杂度完成此题?
// 记录状态
// 第一个数肯定符合
const wiggleMaxLength = (nums) => {
    let res = 1;
    let isUp;
    for (let i = 0; i < nums.length - 1; i++) {
        if (nums[i] === nums[i + 1]) {
            continue;
        }
        else if (nums[i] < nums[i + 1]) {
            if (isUp === undefined)
                (isUp = false), res++;
            isUp && res++, (isUp = !isUp);
        }
        else {
            if (isUp === undefined)
                (isUp = true), res++;
            !isUp && res++, (isUp = !isUp);
        }
    }
    return res;
};
console.log(wiggleMaxLength([1, 7, 4, 9, 2, 5]));
