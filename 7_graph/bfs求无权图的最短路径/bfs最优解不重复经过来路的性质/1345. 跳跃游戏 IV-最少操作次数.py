from collections import defaultdict, deque
from typing import List


class Solution:
    def minJumps(self, arr: List[int]) -> int:
        """请你返回到达数组最后一个元素的下标处所需的 最少操作次数 。
        
        每次可以向左、向右、或者走到值相等的下标处
        注意bfs剪枝
        """
        indexMap = defaultdict(list)
        for i, num in enumerate(arr):
            indexMap[num].append(i)

        n = len(arr)
        queue = deque([(0, 0)])
        visited = set([0])
        while queue:
            cur, step = queue.popleft()
            if cur == n - 1:
                return step

            next1 = cur + 1
            if next1 < n and next1 not in visited:
                visited.add(next1)
                queue.append((next1, step + 1))

            next2 = cur - 1
            if next2 >= 0 and next2 not in visited:
                visited.add(next2)
                queue.append((next2, step + 1))

            for next3 in indexMap[arr[cur]]:
                if next3 not in visited:
                    visited.add(next3)
                    queue.append((next3, step + 1))
            # 剪枝
            # indexMap[arr[cur]].clear()
            indexMap[arr[cur]] = []

        return -1


if __name__ == '__main__':
    # clear是O(n)的
    import time

    for size in 1_000_000, 10_000_000, 100_000_000, 1_000_000_000:
        ls = [None] * size
        t0 = time.time()
        ls.clear()
        t1 = time.time()
        print(size, t1 - t0)

