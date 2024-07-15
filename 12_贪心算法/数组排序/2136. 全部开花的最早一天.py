# 2136. 全部开花的最早一天 (机器A串行，机器B并行/种花问题)
# 播种后，花可以自动生长.
# https://leetcode-cn.com/problems/earliest-possible-day-of-full-bloom/comments/1323899
# https://www.luogu.com.cn/problem/P1248
# 贪心还是要考虑排序
# 1 <= n <= 1e5
# 1 <= plantTime[i], growTime[i] <= 1e4

from typing import List


class Solution:
    def earliestFullBloom(self, plantTime: List[int], growTime: List[int]) -> int:
        n = len(plantTime)
        order = list(range(n))
        order.sort(key=lambda x: growTime[x], reverse=True)  # 按照生长时间从大到小排序
        pre, res = 0, 0
        for i in order:
            a, b = plantTime[i], growTime[i]
            res = max(res, pre + a + b)
            pre += a
        return res


if __name__ == "__main__":
    assert Solution().earliestFullBloom(plantTime=[3, 5, 8, 7, 10], growTime=[6, 2, 1, 4, 9]) == 34
