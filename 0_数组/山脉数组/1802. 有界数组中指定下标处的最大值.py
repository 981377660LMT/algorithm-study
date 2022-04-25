# 你需要构造一个同时满足下述所有条件的数组 nums

# 长为n,元素为正整数，相邻差不超过1，和不超过maxSum,index处取最大值


# 1 <= n <= maxSum <= 10^9
# 返回index处的最大值

# 总结：
# 1. 先 maxSum-n 此时元素只要>=0即可
# 2. 二分答案+山峰型数组
# 指定下标处的最大值


class Solution:
    def maxValue(self, n: int, index: int, maxSum: int) -> int:
        def check(mid: int) -> bool:
            # 0-index-1:1 1 1 1 mid-2 mid-1
            leftCount = min(index, mid - 1)
            leftRemain = index - leftCount
            leftSum = (mid - 1 + mid - leftCount) * leftCount // 2 + leftRemain
            rightCount = min(n - index - 1, mid - 1)
            rightRemain = n - index - 1 - rightCount
            rightSum = (mid - 1 + mid - rightCount) * rightCount // 2 + rightRemain
            return (leftSum + rightSum + mid) <= maxSum

        left, right = 1, maxSum
        while left <= right:
            mid = (left + right) >> 1
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1
        return right


print(Solution().maxValue(n=4, index=2, maxSum=6))
# 输出：2
# 解释：数组 [1,1,2,1] 和 [1,2,2,1] 满足所有条件。不存在其他在指定下标处具有更大值的有效数组。
print(Solution().maxValue(n=4, index=0, maxSum=4))
