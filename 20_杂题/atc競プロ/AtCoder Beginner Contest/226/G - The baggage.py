"""货物搬运分配问题"""

# 有一些货物,重量1-5 和一些人,载重1-5
# 每个人可以搬一些包裹,但是总重量不能大于载重
# 问所有人能否搬完所有包裹
# T<=5e4 nums[i]<=1e16

# !贪心
# https://www.cnblogs.com/7KByte/p/15535262.html
import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def solve() -> bool:
    def assign(weight: int, strength: int) -> None:
        """将重量为weight的包裹分配给力量为strength的人"""
        cur = min(parcels[weight], people[strength])
        parcels[weight] -= cur
        people[strength] -= cur
        people[strength - weight] += cur

    parcels = [0] + list(map(int, input().split()))  # 重量 1, 2, 3, 4, 5 的包裹数量
    people = [0] + list(map(int, input().split()))  # 力量 1, 2, 3, 4, 5 的人数

    assign(5, 5)

    assign(4, 4)  # 如果没有 b4，那么我们需要且必须使用 b5，同时在 b1 中加上剩余空间。
    assign(4, 5)

    assign(3, 3)
    assign(3, 5)  # 先考虑力量为5的,如果使用 b4​，我们还得考虑是否会有 2×a2→a4​ 的情况
    assign(3, 4)

    for i in range(5, 1, -1):
        assign(2, i)
    for i in range(5, 0, -1):
        assign(1, i)

    return all(x == 0 for x in parcels)


if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        print("Yes" if solve() else "No")
