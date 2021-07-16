"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// 有长度为n 的数组，其元素都是int型整数（有正有负）。在连续的子数组中找到其和为最大值的数组。
// 给出数组[-2,2,−3,4,−1,2,1,−5,3]，符合要求的子数组为[4,−1,2,1]，其最大和为 6
// 分析：
// 1.最大子数组的第一个元素、最后一个元素肯定为正值。
// 2.一旦current_sum<0是不是current_sum中记录的连续元素就绝对不可能成为最大子数组，
// 所以此时重置current_sum=0开始从下一个元素记录其连续片段的和。
const maxSubArray = (nums) => {
    if (nums.length <= 1)
        return nums;
    //left_pos,right_pos 记录最大子数组的位置信息
    let left_pos = 0;
    let right_pos = 0;
    let maxSum = 0;
    let curSum = 0;
    for (let index = 1; index < nums.length; index++) {
        curSum += nums[index];
        if (curSum > maxSum) {
            maxSum = curSum;
            right_pos = index;
        }
        if (curSum < 0) {
            curSum = 0;
            left_pos = index + 1;
        }
    }
    return nums.slice(left_pos, right_pos + 1);
};
console.log(maxSubArray([-2, 2, -3, 4, -1, 2, 1, -5, 3]));
