# 两个序列 问能否构造出一个新序列使得 nums[i] 来自 A 或 B
# 且 abs(nums[i]-nums[i+1])<=k


# 滚动数组优化

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, k = map(int, input().split())
    A = list(map(int, input().split()))
    B = list(map(int, input().split()))

    ok1, ok2 = True, True
    for i in range(1, n):
        ok1, ok2 = (
            (ok1 and abs(A[i] - A[i - 1]) <= k) or (ok2 and abs(A[i] - B[i - 1]) <= k),
            (ok1 and abs(B[i] - A[i - 1]) <= k) or (ok2 and abs(B[i] - B[i - 1]) <= k),
        )

    print(("No", "Yes")[ok1 or ok2])


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
