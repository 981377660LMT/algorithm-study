# 将x`未满`的视为小孩，将x`以上`的视为大人
# 任取x，求判断对的最大数量
# !注意，日语里超過/未满　指的是 more/less than
# !以上/以下指的是  more/less than or equal to

# !排序后 枚举判断正确的大人个数 二分出小孩的正确个数 注意判断边界
# !二分统计也可以用前缀和替代
from bisect import bisect_left
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


def main() -> None:
    n = int(input())
    s = input()
    weights = list(map(int, input().split()))
    children, adults = [], []
    for flag, weight in zip(s, weights):
        if flag == "0":
            children.append(weight)
        else:
            adults.append(weight)
    children.sort(), adults.sort()

    if not children or not adults:  # !注意判断边界
        print(n)
        exit(0)

    res = 0
    for pivot in weights:
        ok1 = len(adults) - bisect_left(adults, pivot)
        ok2 = bisect_left(children, pivot)
        res = max(res, ok1 + ok2)

    print(res)


while True:
    try:
        main()
    except (EOFError, ValueError):
        break
