from typing import List

# 1 <= time.length <= 60000
class Solution:
    def numPairsDivisibleBy60(self, time: List[int]) -> int:
        counter = [0] * 60
        res = 0
        for t in time:
            res += counter[-(t % 60)]
            counter[t % 60] += 1

        return res


print(Solution().numPairsDivisibleBy60([30, 20, 150, 100, 40]))
# 输出：3
# 解释：这三对的总持续时间可被 60 整数：
# (time[0] = 30, time[2] = 150): 总持续时间 180
# (time[1] = 20, time[3] = 100): 总持续时间 120
# (time[1] = 20, time[4] = 40): 总持续时间 60

