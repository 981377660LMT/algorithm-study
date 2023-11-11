# 给定长度为N的序列A = (Ao,A1,. . . ,AN-1)。
# 有一个空盘子。Takahashi每次会在其中加入A(x mod N)颗糖果(X是当前盘子中糖果的数量)。
# 求K次操作后的糖果总数。注意糖果总数不要模N
# 2 <= N  ≤= 2e5
# 1 ≤= K  ≤= 1e12
# 1 ≤= Ai ≤= 1e6


# !鸽巢原理 因为是模n 所以n后肯定会进入某个周期循环
# !n天后的牢房

import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, k = map(int, input().split())
    nums = list(map(int, input().split()))

    remain = k
    res, indexMod = 0, 0
    visited = dict()  # 索引 => (糖果数量, 剩余轮数)
    while remain:
        visited[indexMod] = (res, remain)  # !保存当前状态

        res += nums[indexMod]  # !线性转移
        remain -= 1
        indexMod = res % n

        if indexMod in visited:  # !寻找周期加速
            preRes, preRemain = visited[indexMod]
            period = preRemain - remain
            div, mod = divmod(remain, period)
            res += (res - preRes) * div
            remain = mod

    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
