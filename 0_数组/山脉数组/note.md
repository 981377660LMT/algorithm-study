1. 一般题:left right 二分查找即可;注意边界条件
   `162. 寻找峰值.ts`：导数题,你可以假设 nums[-1] = nums[n] = -∞ 。
   `852. 山脉数组的峰顶索引.js`：模板即可
2. 难一点的：`枚举山顶`，动态规划,`统计 up 和 down`
   `845. 数组中的最长山脉.py`
   We have already many 2-pass or 3-pass problems, like `821. Shortest Distance to a Character.`
   They have almost the same idea.
   One forward pass and one backward pass.
   Maybe another pass to get the final result, or you can merge it in one previous pass.

3. `1095. 山脉数组中查找目标值.py`
   先二分找到山顶 再二分查找需要的 target
