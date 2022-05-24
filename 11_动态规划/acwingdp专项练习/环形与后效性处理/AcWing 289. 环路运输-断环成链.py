# 破环成链 + 单调队列优化
# AcWing 289. 环路运输
# 对环上两点x,y 求score[x]+score[y]+distOncycle(x,y) 最大值
# distOncycle(i,j)=min(|i−j|,N−|i−j|)
# n<=1e6

from collections import deque
from typing import Any, List

# 枚举每个点只需枚举一边，因为(i,j)和(j,i)是一样的
# 在i前winLen的一段内找score[i]+score[j]+i-j的最大值
# 即找 score[j]-j 的最大值 滑窗维护最值即可 单调队列
def calMax1(scores: List[int]) -> int:
    """deque写法"""
    n = len(scores)
    res, winLen = 0, n // 2
    scores = [*scores, *scores]  # 破环成链

    maxQueue = deque()
    for i in range(2 * n):
        while maxQueue and maxQueue[0][1] < i - winLen:
            maxQueue.popleft()
        if maxQueue:
            res = max(res, scores[i] + i + maxQueue[0][0])
        while maxQueue and maxQueue[-1][0] < (scores[i] - i):
            maxQueue.pop()
        maxQueue.append([scores[i] - i, i])
    return res


def calMax2(scores: List[int]) -> int:
    """maxQueue写法"""
    n = len(scores)
    res, winLen = 0, n // 2
    scores = [*scores, *scores]  # 破环成链

    maxQueue = MaxQueue()
    for i in range(2 * n):
        while maxQueue and maxQueue[0][1] < i - winLen:
            maxQueue.popleft()
        if maxQueue:
            res = max(res, scores[i] + i + maxQueue.max)
        maxQueue.append(scores[i] - i, i)
    return res


n = int(input())
scores = list(map(int, input().split()))
print(calMax2(scores))


class MaxQueue(deque):
    @property
    def max(self) -> int:
        if not self:
            raise ValueError('maxQueue is empty')
        return self[0][0]

    def append(self, value: int, *metaInfo: Any) -> None:
        count = 1
        while self and self[-1][0] < value:
            count += self.pop()[-1]
        super().append([value, *metaInfo, count])

    def popleft(self) -> None:
        if not self:
            raise IndexError('popleft from empty queue')

        self[0][-1] -= 1
        if self[0][-1] == 0:
            super().popleft()

