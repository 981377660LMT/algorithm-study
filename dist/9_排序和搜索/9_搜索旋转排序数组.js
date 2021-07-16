"use strict";
// 给你 旋转后 的数组 nums 和一个整数 target ，如果 nums 中存在这个目标值 target ，则返回它的下标，否则返回 -1 。
Object.defineProperty(exports, "__esModule", { value: true });
// 注意到数组可能是有序的:二分查找log(n)复杂度
// 需要在一段有序的数组上使用二分查找 注意必须是有序的
// 这样就知道target在左边还是右边
// 注意是>=而不是>
const search = (nums, target) => {
    let left = 0;
    let right = nums.length - 1;
    while (left <= right) {
        const mid = Math.floor((left + right) / 2);
        if (nums[mid] === target)
            return mid;
        // if (nums[left] === target) return nums[left]
        // if (nums[right] === target) return nums[right]
        // 左半部分有序
        if (nums[left] <= nums[mid]) {
            if (nums[mid] >= target && target >= nums[left]) {
                right = mid - 1;
            }
            else {
                // target不在左半部分
                left = mid + 1;
            }
        }
        else {
            // 右半部分有序
            if (nums[right] >= target && target >= nums[mid]) {
                left = mid + 1;
            }
            else {
                // target不在右半部分
                right = mid - 1;
            }
        }
    }
    return -1;
};
console.log(search([4, 5, 6, 7, 0, 1, 2], 0));
