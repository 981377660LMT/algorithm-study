# 给你一个下标从 0 开始的整数数组 costs ，其中 costs[i] 是雇佣第 i 位工人的代价。

# 同时给你两个整数 k 和 candidates 。我们想根据以下规则恰好雇佣 k 位工人：

# 总共进行 k 轮雇佣，且每一轮恰好雇佣一位工人。
# 在每一轮雇佣中，从最前面 candidates 和最后面 candidates 人中选出代价最小的一位工人，如果有多位代价相同且最小的工人，选择下标更小的一位工人。
# 比方说，costs = [3,2,7,7,1,2] 且 candidates = 2 ，第一轮雇佣中，我们选择第 4 位工人，因为他的代价最小 [3,2,7,7,1,2] 。
# 第二轮雇佣，我们选择第 1 位工人，因为他们的代价与第 4 位工人一样都是最小代价，而且下标更小，[3,2,7,7,2] 。注意每一轮雇佣后，剩余工人的下标可能会发生变化。
# 如果剩余员工数目不足 candidates 人，那么下一轮雇佣他们中代价最小的一人，如果有多位代价相同且最小的工人，选择下标更小的一位工人。
# 一位工人只能被选择一次。
# 返回雇佣恰好 k 位工人的总代价。


from typing import List
from sortedcontainers import SortedList


class Solution:
    def totalCost(self, costs: List[int], k: int, candidates: int) -> int:
        people = [(costs[i], i) for i in range(len(costs))]
        sl = SortedList()
        if candidates * 2 >= len(costs):
            sl = SortedList(people)
        else:
            sl.update(people[:candidates])
            sl.update(people[-candidates:])

        # 头尾指针表示当前看到哪个位置的人
        left, right = candidates - 1, len(costs) - candidates
        res = 0
        for _ in range(k):
            min_, index = sl.pop(0)
            res += min_
            if index <= left:
                left += 1
                if left < right:
                    sl.add(people[left])
            else:
                right -= 1
                if left < right:
                    sl.add(people[right])
        return res
