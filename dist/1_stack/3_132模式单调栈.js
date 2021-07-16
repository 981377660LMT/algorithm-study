"use strict";
// 给你一个整数数组 nums ，数组中共有 n 个整数。
// 132 模式的子序列 由三个整数 nums[i]、nums[j] 和 nums[k] 组成，并同时满足：
// i < j < k 和 nums[i] < nums[k] < nums[j] 。
Object.defineProperty(exports, "__esModule", { value: true });
// 定一找二(找32形式的)
const find132Pattern = (nums) => {
    if (nums.length <= 2)
        return false;
    // 存储k,k是一个单调递减的栈
    const stack = [];
    // 存储每个位置左侧的最小值i
    const minI = [nums[0]];
    for (let index = 1; index < nums.length; index++) {
        const element = nums[index];
        minI.push(Math.min(minI[minI.length - 1], element));
    }
    console.log(minI);
    for (let j = nums.length - 1; j >= 1; j--) {
        if (nums[j] > minI[j]) {
            // 去除不合适的k
            while (stack.length && stack[stack.length - 1] <= minI[j])
                stack.pop();
            // 找到了合适的k
            if (stack.length && stack[stack.length - 1] < nums[j])
                return true;
            stack.push(nums[j]);
        }
    }
    return false;
};
console.log(find132Pattern([3, 5, 0, 3, 4]));
console.log(find132Pattern([3, 1, 4, 2]));
