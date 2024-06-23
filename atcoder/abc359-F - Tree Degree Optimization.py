# abc359-F - Tree Degree Optimization (优先队列贪心)
# 给定n个点的点权nums[i],构造一棵树，使得 ∑(deg[i]^2) * nums[i] 最小
# n<=2e5, nums[i]<=1e9.
#
# 注意到树的度数: 1 <= deg[i] <= n-1, 且 ∑deg[i] = 2*(n-1)
# !可以把问题抽象成每个点起始度数为1，把剩下的(n-2)个度分配给每个点，使得代价最小
# !每次仅分配1的度，那肯定是贪心的分配给代价最小的点。
# !优先队列维护(每个点度数+1后增长的代价)

from heapq import heapify, heappop, heappush
from typing import List


def treeDegreeOptimization(nums: List[int]) -> int:
    n = len(nums)
    deg = [1] * n
    pq = [(3 * v, i) for i, v in enumerate(nums)]  # (增长代价, 点编号)
    heapify(pq)
    for _ in range(n - 2):
        _, id = heappop(pq)
        deg[id] += 1
        d = deg[id]
        delta = ((d + 1) * (d + 1) - d * d) * nums[id]
        heappush(pq, (delta, id))
    return sum(d * d * v for d, v in zip(deg, nums))


if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    print(treeDegreeOptimization(nums))
