# 寻找循环的起点
# !注意这里要先进入循环轨道(即第一次走回了看过的点) 再找周期
# 1e100个土豆(无限个土豆)的土豆流 n个一组重量循环(n<=2e5)
# 打包土豆 如果一组土豆重量>=x 那么就打包到下一组
# Q个询问(Q<=2e5) 求第k组土豆的数量 (1<=k<=1e12)

# !1. 滑窗记录每个土豆作为左端点，最右边能到哪个土豆 (每个土豆开始的组能放几个土豆;滑窗处理环比较方便)
# !2. 进入循环轨道后哈希表找周期(当然也可以倍增dp)
from itertools import accumulate
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, q, limit = map(int, input().split())
    weights = [int(num) for num in input().split()]

    # # 滑窗查找每个土豆开始的组能放几个土豆
    div, mod = divmod(limit, sum(weights))  # 注意先要模
    size = [div * n] * n
    right, curSum = 0, 0
    for left in range(n):
        while curSum + (weights[right % n]) < mod:
            curSum += weights[right % n]
            right += 1
        size[left] += right - left + 1
        curSum -= weights[left]

    left, visited = 0, dict()  # 字典避免线性查找
    queryLeft = []  # 每个询问的左端点土豆
    while True:
        if left in visited:
            break
        visited[left] = len(visited)
        queryLeft.append(left)
        left = (left + size[left]) % n
    start = visited[left]  # 循环起点索引
    freq = len(queryLeft) - start  # 循环周期

    for _ in range(q):
        k = int(input()) - 1
        if k >= start:
            k = ((k - start) % freq) + start
        left = queryLeft[k]
        print(size[left])


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
