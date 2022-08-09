# 苺个任务都可以在 1 个单位时间内执行完
# 两个 相同种类 的任务之间必须有长度为整数 k 的冷却时间，
# 因此至少有连续 k 个单位时间内 CPU 在执行不同的任务，或者在待命状态。
# 计算完成所有任务所需要的 最短时间


# 假设有无数个桶
# 我们设计桶的大小为 n+1，则相同的任务恰好不能放入同一个桶，最密也只能放入相邻的桶。
# 一个桶不管是否放满，其占用的时间均为 n+1，这是因为后面桶里的任务需要等待冷却时间。
# 最后一个桶是个特例，由于其后没有其他任务需等待，所以占用的时间为桶中的任务个数。
# !总排队时间 = (桶个数 - 1) * (n + 1) + 最后一桶的任务数
# 如果任务太多桶放不下，则取tasks.length

from collections import Counter
from typing import List


class Solution:
    def leastInterval(self, tasks: List[str], k: int) -> int:
        """分桶  完成所有任务的最短时间取决于出现次数最多的任务数量。"""
        counter = Counter(tasks)
        freqCounter = Counter(counter.values())
        max_ = max(freqCounter)  # 最多的任务数量
        maxCount = freqCounter[max_]  # 最多的任务数量的个数
        res = (max_ - 1) * (k + 1) + maxCount
        return max(res, len(tasks))


print(Solution().leastInterval(tasks=["A", "A", "A", "B", "B", "B"], k=2))
