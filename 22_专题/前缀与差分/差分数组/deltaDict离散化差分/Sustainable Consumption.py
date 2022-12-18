# 问生产者是否可以供应消费者的需求
import collections


class Solution:
    def solve(self, producers, consumers):
        deltaDict = collections.defaultdict(int)
        for s, e, j in producers:
            deltaDict[s] += j
            deltaDict[e + 1] -= j
        for s, e, j in consumers:
            deltaDict[s] -= j
            deltaDict[e + 1] += j

        preSum = 0
        for key in sorted(deltaDict):
            preSum += deltaDict[key]
            if preSum < 0:
                return False
        return True


print(Solution().solve(producers=[[0, 10, 5], [5, 15, 10]], consumers=[[5, 10, 15], [11, 15, 8]]))
