# 你需要构造一个同时满足下述所有条件的数组 nums

# 长为n,元素为正整数，相邻差不超过1，和不超过maxSum,index处取最大值
# 1 <= n <= maxSum <= 10^9
# 返回index处的最大值

# 总结：
# 1. 先 maxSum-n 此时元素只要>=0即可
# 2. 二分答案+山峰型数组
class Solution:
    def maxValue(self, n: int, index: int, maxSum: int) -> int:
        def check(maxVal: int) -> bool:
            leftMin = max(maxVal - index, 0)
            leftSum = (leftMin + maxVal) * (maxVal - leftMin + 1) / 2
            rightMin = max(maxVal - ((n - 1) - index), 0)
            rightSum = (maxVal + rightMin) * (maxVal - rightMin + 1) / 2
            return (leftSum + rightSum - maxVal) <= maxSum

        maxSum -= n
        left, right = 0, maxSum
        while left <= right:
            mid = (left + right) >> 1
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1
        return left


print(Solution().maxValue(n=4, index=2, maxSum=6))
# 输出：2
# 解释：数组 [1,1,2,1] 和 [1,2,2,1] 满足所有条件。不存在其他在指定下标处具有更大值的有效数组。
