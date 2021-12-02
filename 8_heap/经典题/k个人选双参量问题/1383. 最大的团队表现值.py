from typing import List
from heapq import heappop, heappush

# 请你返回由最多 k 个工程师组成的 ​​​​​​最大团队表现值 ，由于答案可能很大，请你返回结果对 10^9 + 7 取余后的结果。
# 团队表现值 的定义为：一个团队中「所有工程师速度的和」乘以他们「效率值中的最小值」。
# 1 <= n <= 10^5

MOD = int(1e9 + 7)


# 2个变量，想办法”定“住一个====提前把一个排好序(涉及到最值的变量) 取的时候已经是最小了
# 总结:乘参量1用排序，加参量2用堆维护,出堆入堆形成抗衡，同时更新res
# 题目中 对eff降序排列，堆speed升序维护:每次入堆eff最没用的，出堆speed最没用的,出堆入堆形成抗衡
class Solution:
    def maxPerformance(self, n: int, speed: List[int], efficiency: List[int], k: int) -> int:
        pq = []
        res = spdSum = 0

        for eff, spd in sorted(zip(efficiency, speed), reverse=True):
            heappush(pq, spd)
            spdSum += spd
            if len(pq) > k:
                spdSum -= heappop(pq)
            res = max(res, spdSum * eff)

        return res % MOD


print(Solution().maxPerformance(n=6, speed=[2, 10, 3, 1, 5, 8], efficiency=[5, 4, 3, 9, 7, 2], k=2))
# 输入：n = 6, speed = [2,10,3,1,5,8], efficiency = [5,4,3,9,7,2], k = 2
# 输出：60
# 解释：
# 我们选择工程师 2（speed=10 且 efficiency=4）和工程师 5（speed=5 且 efficiency=7）。他们的团队表现值为 performance = (10 + 5) * min(4, 7) = 60 。
