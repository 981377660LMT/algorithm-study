# 在游戏过程中，可以把 x 与其上、下、左、右四个方向之一的数字交换（如果存在）。
# 我们的目的是通过交换，使得网格变为如下排列（称为正确排列）：
# 例如：

# 1 2 3
# x 4 6
# 7 5 8
# 交换到：
# 1 2 3
# 4 5 6
# 7 8 x

# 用字符串作为状态比较方便
from collections import deque

dx = [0, 1, 0, -1]
dy = [1, 0, -1, 0]
target = "12345678x"


def bfs(start: str) -> int:
    visited = set([start])
    queue = deque([(start, 0)])
    while queue:
        cur, step = queue.popleft()
        if cur == target:
            return step
        index = cur.find('x')
        x, y = divmod(index, 3)

        chars = list(cur)
        for i in range(4):
            nx, ny = x + dx[i], y + dy[i]
            if 0 <= nx < 3 and 0 <= ny < 3:
                nIndex = nx * 3 + ny
                chars[nIndex], chars[index] = chars[index], chars[nIndex]
                next = "".join(chars)
                if next not in visited:
                    visited.add(next)
                    queue.append((next, step + 1))
                # 还原现场
                chars[nIndex], chars[index] = chars[index], chars[nIndex]
    return -1


start = ''.join(list(input().split()))
print(bfs(start))
