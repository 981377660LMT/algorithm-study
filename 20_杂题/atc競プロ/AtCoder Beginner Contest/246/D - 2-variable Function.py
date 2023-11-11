# 要求大于等于n(n≤1e18) 的最小的可以分解成 x = a^3 + a^2*b + a*b^2 + b^3 的数 x (其中a,b为非负整数).

# 首先根据n的范围我们可以确定a和b一定都小于1e6,
# 所以我们可以直接枚举a的取值,
# 用二分法寻找最小的满足的要求的b,然后所有结果里取最小值就行了.


import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)

UPPER = int(1e6 + 5)


def cal(a: int, b: int) -> int:
    return a * a * a + a * b * b + a * a * b + b * b * b


def main() -> None:
    n = int(input())
    res = int(1e18)

    for i in range(UPPER):
        left, right = 0, UPPER
        while left <= right:
            mid = (left + right) // 2
            if cal(i, mid) < n:
                left = mid + 1
            else:
                right = mid - 1

        cand = cal(i, left)
        if cand >= n:
            res = min(res, cand)

    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
