from collections import deque
from typing import Deque, List, Tuple

# 在编号为 i 弹簧处按动弹簧，小球可以弹向 0 到 i-1
# 中任意弹簧或者 i+jump[i] 的弹簧（若 i+jump[i]>=N ，则表示小球弹出了机器）。
# 小球位于编号 0 处的弹簧时不能再向左弹。

# 你需要将小球弹出机器。请求出最少需要按动多少次弹簧，可以将小球从编号 0 弹簧弹出整个机器，即向右越过编号 N-1 的弹簧。

# 由于我们是 BFS，一个位置只要被更新了，就一定不需要再次被更新。
# 所以，我们维护一个变量 preidx，表示当前 [1, preidx - 1]之间的所有位置都已经被更新过了。
# 这样，我们之后只需要扩展 preidx 右边，idx 左边的位置即可。
# 这样最坏情况下只会遍历全部弹簧位置，复杂度可以保持在O(n)。


class Solution:
    def minJump(self, jump: List[int]) -> int:
        """请求出最少需要按动多少次弹簧，可以将小球从编号 0 弹簧弹出整个机器
        
        bfs求最短路+剪枝
        """
        n = len(jump)
        queue: Deque[Tuple[int, int]] = deque([(0, 0)])
        visited = set([0])
        # 维护一个变量 preidx，表示当前 [1, preidx - 1]之间的所有位置都已经被更新过了。
        preMax = 0

        while queue:
            cur, curDist = queue.popleft()

            # 右跳jump[i]
            next1 = cur + jump[cur]
            if next1 >= n:
                return curDist + 1
            else:
                if next1 not in visited:
                    queue.append((next1, curDist + 1))
                    visited.add(next1)

            # 左跳任意长度
            for next2 in range(preMax, cur):
                if next2 not in visited:
                    queue.append((next2, curDist + 1))
                    visited.add(next2)

            preMax = max(cur, preMax)

        return -1
