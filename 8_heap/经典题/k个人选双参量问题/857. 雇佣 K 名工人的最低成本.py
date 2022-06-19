from typing import List
from heapq import heappop, heappush

# 有 N 名工人。 第 i 名工人的工作质量为 quality[i] ，其最低期望工资为 wage[i] 。
# 1 <= K <= N <= 10000
# 对工资组中的每名工人，应当按其工作质量与同组其他工人的工作质量的比例来支付工资。
# 工资组中的每名工人至少应当得到他们的最低期望工资。
# 返回组成一个满足上述条件的工资组至少需要多少钱。
# https://leetcode-cn.com/problems/minimum-cost-to-hire-k-workers/comments/1031682


# 我们发工资肯定希望wage/quality 计件单价 越小越好
# 总结:乘参量1用排序，加参量2用堆维护,出堆入堆形成抗衡，同时更新res
# 题目中 对wage/quality升序排列，对quality降序维护:每次入堆wage/quality变大，出堆quality变小,出堆入堆形成抗衡


class Solution:
    def mincostToHireWorkers(self, quality: List[int], wage: List[int], k: int) -> float:
        workers = sorted((float(w) / q, q) for w, q in zip(wage, quality))
        res = int(1e20)
        quaSum = 0
        pq = []
        for wa, qua in workers:
            heappush(pq, -qua)
            quaSum += qua
            if len(pq) > k:
                quaSum -= -heappop(pq)
            if len(pq) == k:
                res = min(res, quaSum * wa)
        return res


print(Solution().mincostToHireWorkers(quality=[10, 20, 5], wage=[70, 50, 30], k=2))
# 输出： 105.00000
# 解释： 我们向 0 号工人支付 70，向 2 号工人支付 35。


# quality应该翻译成工作量/搬砖数；
# 工资标准为统一的计件单价的计件工资，同工同酬。每个人的计件结果，不能低于期望工资，否则候选人不接Offer；
# 有k个工位，老板一定要看到坐满人才开心，要不会认为HR能力低下、招不到人；
# 产出搁一边，那是生产单位的问题，反正老板要求工资总额支出一定要少

# 先把候选人按计件单价排序，从低的开始捋。
# 在每一个单价的范围内，找出K个搬砖最少的，这样工资总额比较少。
