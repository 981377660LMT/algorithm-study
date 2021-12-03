from typing import List
from heapq import heapify, heappop, heappush

# blocks[i] = t 意味着第  i 个街区需要 t 个单位的时间来建造。
# 一个工人要么需要再召唤一个工人（工人数增加 1）；要么建造完一个街区后回家。这两个决定都需要花费一定的时间。
# 一个工人再召唤一个工人所花费的时间由整数 split 给出。
# 如果两个工人同时召唤别的工人，那么他们的行为是并行的，所以时间花费仍然是 split。
# 最开始的时候只有 一个 工人，请你最后输出建造完所有街区所需要的最少时间。

# 总结：
# 哈夫曼树(贪心思想):为了让数值大的节点尽可能少参与到合并中，我们总是优先挑选两个最小的节点来进行合并。
# 每个工人可以看作一个树的节点，可以分裂，变成两个节点
# 每次选取两个最小的点，合并成一个新的点
# 最后root的值就是总的最小值


class Solution:
    def minBuildTime(self, blocks: List[int], split: int) -> int:
        pq = blocks[:]
        heapify(pq)
        while len(pq) > 1:
            a, b = heappop(pq), heappop(pq)
            heappush(pq, max(a, b) + split)
        return pq[0]


print(Solution().minBuildTime(blocks=[1, 2], split=5))
# 输出：7
# 解释：我们用 5 个时间单位将这个工人分裂为 2 个工人，然后指派每个工人分别去建造街区，从而时间花费为 5 + max(1, 2) = 7
print(Solution().minBuildTime(blocks=[1, 2, 3], split=1))
# 输出：4
# 解释：
# 将 1 个工人分裂为 2 个工人，然后指派第一个工人去建造最后一个街区，并将第二个工人分裂为 2 个工人。
# 然后，用这两个未分派的工人分别去建造前两个街区。
# 时间花费为 1 + max(3, 1 + max(1, 2)) = 4
