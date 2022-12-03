from itertools import product
import sys
from typing import List, Tuple

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")

# 从团队中选出3个人,这个队伍每门科目的的得分都是这3个人中最高的那个人的得分
# 将x记为五门课的得分的最小值,问x最大是多少

# !最大化最小值 => 二分答案
# !固定mid后 每个人的能力为0/1 `能力数只有32种情况` 枚举所有情况即可


def madTeam(grid: List[Tuple[int, int, int, int, int]]) -> int:
    def check(mid: int) -> bool:
        states = set()
        for person in grid:
            states.add(sum(1 << i for i in range(5) if person[i] >= mid))
        return any((a | b | c == TARGET) for (a, b, c) in product(states, repeat=3))

    TARGET = (1 << len(grid[0])) - 1  # 5个科目都ok
    left, right = 1, int(1e9 + 10)
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            left = mid + 1
        else:
            right = mid - 1

    return right


if __name__ == "__main__":
    n = int(input())
    grid = [tuple(map(int, input().split())) for _ in range(n)]
    print(madTeam(grid))
