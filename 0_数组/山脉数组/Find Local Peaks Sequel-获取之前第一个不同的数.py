INF = int(1e20)

# 只要他比左右两个第一个不同的数都大，那么他就是peak
class Solution:
    def solve(self, nums):
        def getFirstDiff(nums):
            """获取之前第一个不同的数"""
            res = []
            pre = -INF
            for num in nums:
                if num != pre:
                    res.append(pre)
                    pre = num
                else:
                    res.append(res[-1])

            return res

        n = len(nums)
        left = getFirstDiff(nums)
        right = getFirstDiff(nums[::-1])[::-1]

        res = []

        for i in range(n):
            if left[i] == right[i] == -INF:
                continue
            if left[i] < nums[i] > right[i]:
                res.append(i)

        return res
