# 请你能返回设置 targetSeconds 秒钟加热时间需要花费的最少代价。
# bfs暴力
from heapq import heappop, heappush

# todo
class Solution:
    def minCostSetTime(self, startAt: int, moveCost: int, pushCost: int, targetSeconds: int) -> int:
        def cal(path: list[int]) -> int:
            # print(path)
            digits = ''.join([str(num) for num in path]).zfill(4)
            if -1 in path:
                return -0x3F3F3F3F
            minute, second = int(digits[:2]), int(digits[2:])
            return minute * 60 + second

        # cost, path, index, hasPushed
        queue = [(0, [startAt, -1, -1, -1], 0, False)]
        while queue:
            cost, path, index, hasPushed = heappop(queue)
            if index == 3 and cal(path) == targetSeconds:
                print(path)
                return cost
            if index >= 4:
                continue

            for next in range(index + 1, 4):
                heappush(queue, (cost + moveCost, path, next, False))
            if not hasPushed:
                for select in range(10):
                    path = path[:]
                    path[index] = select
                    heappush(queue, (cost + pushCost, path, index, True))
        return -1


print(Solution().minCostSetTime(startAt=1, moveCost=2, pushCost=1, targetSeconds=600))
print(Solution().minCostSetTime(startAt=0, moveCost=1, pushCost=2, targetSeconds=76))
