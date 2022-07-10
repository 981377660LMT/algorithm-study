# 每次操作会交换
# x元素所在位置与其右边的元素的位置
# 如果右边没有数 与左边交换

# !利用元素各异的性质 哈希表

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, q = map(int, input().split())
    nums = list(range(1, n + 1))
    indexMap = {v: i for i, v in enumerate(nums)}

    for _ in range(q):
        x = int(input())
        i1 = indexMap[x]
        i2 = i1 + 1 if i1 < n - 1 else i1 - 1
        nums[i1], nums[i2] = nums[i2], nums[i1]
        indexMap[nums[i1]], indexMap[nums[i2]] = i1, i2
    print(*nums)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
