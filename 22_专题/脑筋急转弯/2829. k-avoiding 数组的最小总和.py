# 2829. k-avoiding 数组的最小总和
# https://leetcode.cn/problems/determine-the-minimum-sum-of-a-k-avoiding-array/description/

# 给你两个整数 n 和 k 。
# !对于一个由 不同 正整数组成的数组，如果其中不存在任何求和等于 k 的不同元素对，则称其为 k-avoiding 数组。
# 返回长度为 n 的 k-avoiding 数组的可能的最小总和。


class Solution:
    def minimumSum2(self, n: int, k: int) -> int:
        res = [1]
        while len(res) < n:
            cur = res[-1] + 1
            while any(cur + x == k for x in res):
                cur += 1
            res.append(cur)
        return sum(res)

    def minimumSum(self, n: int, k: int) -> int:
        """
        O(1)数学.
        相加等于 k 的正整数对有 (1, k - 1), (2, k - 2), ...,共 k // 2 对.
        为了让和最小,我们让每个数都尽可能小,即 1,2,3,...,k // 2.
        如果这些数不够,我们再加上 k, k+1, ..., 这 (n - k // 2)个数.
        """

        def getSum(first: int, diff: int, count: int) -> int:
            last = first + (count - 1) * diff
            return (first + last) * count // 2

        half = k // 2
        if half >= n:
            return getSum(1, 1, n)

        return getSum(1, 1, half) + getSum(k, 1, n - half)
