from collections import defaultdict
from heapq import heappop, heappush


n, startMoney = map(int, input().split())
nums = list(map(int, input().split()))
nums = nums[:-1]
lastPrice = nums[-1]

dist = defaultdict(lambda: -int(1e20))
dist[0] = startMoney
pq = [(-startMoney, 0, startMoney, 0)]  # 现在的钱，股票数，index

# dist[i][j]表示，在第i天持有j股的条件下，剩余的最大现金

while pq:
    score, count, money, index = heappop(pq)
    score *= -1
    if index == n - 1:
        print(score)
        exit()

    # 不买
    if score > dist[index + 1]:
        heappush(pq, (-score, count, money, index + 1))
        dist[index + 1] = score

    # 买
    curCost = nums[index]
    if money >= curCost:
        distCand = score + lastPrice - curCost
        if distCand > dist[index + 1]:
            heappush(pq, (-distCand, count + 1, money - curCost, index + 1))
            dist[index + 1] = score

    # 卖
    if count > 0:
        distCand = score + lastPrice
        if distCand > dist[index + 1]:
            heappush(pq, (-distCand, count - 1, money + lastPrice, index + 1))
            dist[index + 1] = score

