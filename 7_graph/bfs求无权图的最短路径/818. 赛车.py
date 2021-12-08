# 1 <= target（目标位置） <= 10000。
# 你的车会根据一系列由 A（加速）和 R（倒车）组成的指令进行自动驾驶 。
# 当车得到指令 "A" 时, 将会做出以下操作： position += speed, speed *= 2。
# 当车得到指令 "R" 时, 将会做出以下操作：如果当前速度是正数，则将车速调整为 speed = -1 ；否则将车速调整为 speed = 1。  (当前所处位置不变。)

# https://leetcode.com/problems/race-car/discuss/762584/Python-3-Simple-Steps-(BFS)
# https://leetcode-cn.com/problems/race-car/solution/cpython3-bfsji-yi-hua-jian-zhi-by-hanxin-veaq/
from collections import deque


# 最短路径bfs

# 加个简单的剪枝
# 因为位置 ≤10000
# -10000 < n_pos < 20000 and abs(n_speed) < 10000


class Solution:
    def racecar(self, target: int) -> int:
        # 位置，速度，花费
        queue = deque([(0, 1, 0)])
        # 位置，速度
        visited = set([(0, 1)])
        while queue:
            pos, speed, cost = queue.popleft()
            if pos == target:
                return cost

            # 加速
            nextPos = pos + speed
            nextSpeed = speed * 2
            if abs(nextSpeed) < 10000:
                if (nextPos, nextSpeed) not in visited:
                    queue.append((nextPos, nextSpeed, cost + 1))
                    visited.add((nextPos, nextSpeed))

            # 倒车
            nextPos = pos
            nextSpeed = -1 if speed > 0 else 1
            if abs(nextSpeed) < 10000:
                if (nextPos, nextSpeed) not in visited:
                    queue.append((nextPos, nextSpeed, cost + 1))
                    visited.add((nextPos, nextSpeed))

        return -1


print(Solution().racecar(target=6))
# 输出: 5
# 解释:
# 最短指令列表为 "AAARA"
# 位置变化为 0->1->3->7->7->6
