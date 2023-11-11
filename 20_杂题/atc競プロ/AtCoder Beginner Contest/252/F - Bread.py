# 给你一个长度为L的面包和一个数组 A ,规定操作“切开”的代价为：
# 将一个长度为k的面包切分成x和k-x的两段，代价为k
# 现在问如何将这个面包切分为数组A，可能会有剩余；
# n<=2e5
# ai<=1e9
# sum(ai)<=L<=1e15

# https://atcoder.jp/contests/abc252/submissions/36006064
# 逆向思维，把分割看成合并。这样就转换成经典的 Huffman 问题，用最小堆实现。
# !如果 sum(a) < L，可以把多余的面包额外当成一个新的小孩需要的。

from heapq import heapify, heappop, heappush
import sys
import os


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    _, L = map(int, input().split())
    nums = list(map(int, input().split()))
    sum_ = sum(nums)
    if sum_ < L:
        nums.append(L - sum_)  # 把多余的面包额外当成一个新的小孩需要的。
    pq = nums[:]
    heapify(pq)

    res = 0
    while len(pq) >= 2:
        x = heappop(pq)
        y = heappop(pq)
        res += x + y
        heappush(pq, x + y)
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
