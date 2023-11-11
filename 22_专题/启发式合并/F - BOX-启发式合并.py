import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")


# 有n(n<=3e5)个箱子和无数个球
# 给出q(q<=3e5)个操作
# 1 x y 把箱子y的球全部倒入箱子x 保证x!=y 1≤x,y≤n
# 2 x 将下一个编号的球放入箱子x 1≤x≤n
# 3 x 输出球x所在的箱子编号
# !启发式合并 swap操作需要维护mp/rmp


def box(n: int, operations: List[List[int]]) -> List[int]:
    boxes = [[i] for i in range(n + 1)]
    mp = {i: i for i in range(1, n + 1)}  # !箱子编号 => 箱子索引
    rmp = {i: i for i in range(1, n + 1)}  # !箱子索引 => 箱子编号
    belong = {i: i for i in range(1, n + 1)}  # !球所在的箱子索引
    cur = n + 1

    res = []
    for kind, *args in operations:
        if kind == 1:
            x, y = args  # !把y箱子的球倒进x箱子
            i1, i2 = mp[x], mp[y]
            if len(boxes[i1]) < len(boxes[i2]):  # !swap
                mp[x], mp[y] = mp[y], mp[x]
                rmp[i1], rmp[i2] = rmp[i2], rmp[i1]
                i1, i2 = i2, i1
                x, y = y, x
            for ball in boxes[i2]:
                belong[ball] = i1
            boxes[i1] += boxes[i2]
            boxes[i2] = []
        elif kind == 2:  # !把cur球放入x箱子
            x = args[0]
            i = mp[x]
            boxes[i].append(cur)
            belong[cur] = i
            cur += 1
        else:  # !输出x球所在的箱子编号
            x = args[0]
            res.append(rmp[belong[x]])

    # print(boxes, belong)
    return res


if __name__ == "__main__":
    n, q = map(int, input().split())
    operations = [list(map(int, input().split())) for _ in range(q)]
    print(*box(n, operations), sep="\n")
