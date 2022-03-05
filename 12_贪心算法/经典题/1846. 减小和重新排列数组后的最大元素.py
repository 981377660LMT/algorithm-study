from typing import List

# arr 中 第一个 元素必须为 1 。
# 任意相邻两个元素的差的绝对值 小于等于 1
# 1 <= arr[i] <= 109

# 你可以执行以下 2 种操作任意次：
# 减小 arr 中任意元素的值，使其变为一个 更小的正整数 。
# 重新排列 arr 中的元素，你可以以任意顺序重新排列。


# 请你返回执行以上操作后，在满足前文所述的条件下，arr 中可能的 最大值 。

# 10^5 考虑贪心/排序/dp/双指针
# 最后肯定是1,2,3,... 每个数要被限制
class Solution:
    def maximumElementAfterDecrementingAndRearranging(self, arr: List[int]) -> int:
        limit = 1
        for val in sorted(arr)[1:]:
            limit = min(limit + 1, val)
        return limit


print(Solution().maximumElementAfterDecrementingAndRearranging(arr=[100, 1, 1000]))
# 输出：3
# 解释：
# 一个可行的方案如下：
# 1. 重新排列 arr 得到 [1,100,1000] 。
# 2. 将第二个元素减小为 2 。
# 3. 将第三个元素减小为 3 。
# 现在 arr = [1,2,3] ，满足所有条件。
# arr 中最大元素为 3 。
