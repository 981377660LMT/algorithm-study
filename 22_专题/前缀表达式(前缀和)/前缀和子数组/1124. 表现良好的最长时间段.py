from typing import List


# 我们认为当员工一天中的工作小时数大于 8 小时的时候，那么这一天就是「劳累的一天」。
# 所谓「表现良好的时间段」，意味在这段时间内，「劳累的天数」是严格 大于「不劳累的天数」。
# 请你返回「表现良好时间段」的`最大长度`。

# 这道题目的本质是 先把数组换成由-1 和1 （-1代表没有超时，1代表超时） 转化为找一个maximum size subarray 让它的和大于 0，
# 525. Subarray Sum Equals K


# 注意不可滑窗 因为没有单向性

# 思路:
# 1. 开一个pre
# 2. 如果cursum>0 更新res
# 3. 如果前缀里有`cursum-1`(贪心) 更新res
# 4. setdefault
class Solution:
    def longestWPI(self, hours: List[int]) -> int:
        pre = dict({0: -1})
        res = cursum = 0

        for i, h in enumerate(hours):
            cursum += 1 if h > 8 else -1
            if cursum > 0:
                res = i + 1
            if cursum - 1 in pre:
                res = max(res, i - pre[cursum - 1])
            pre.setdefault(cursum, i)
        return res


print(Solution().longestWPI([9, 9, 6, 0, 6, 6, 9]))
