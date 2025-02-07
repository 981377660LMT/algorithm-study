# F - K-th Largest Triplet
# https://atcoder.jp/contests/abc391/tasks/abc391_f
# https://atcoder.jp/contests/abc391/editorial/12086
# !给定三个长为N的数组A,B,C,求第K大的三元组A[i]*B[j]+B[j]*C[k]+C[k]*A[i].
# !N<=2e5,1<=K<=min(N,2e5)
# 二分+剪枝

import sys

input = lambda: sys.stdin.readline().rstrip("\r\n")


def kthLargestTriplet(A: list[int], B: list[int], C: list[int], K: int) -> int:
    A, B, C = sorted(A), sorted(B), sorted(C)  # type: ignore
    N = len(A)

    def f(i: int, j: int, k: int) -> int:
        return A[i] * B[j] + B[j] * C[k] + C[k] * A[i]

    def check(mid: int) -> bool:
        """不超过mid的数<k个."""
        count = 0
        for i in range(N):
            if f(i, 0, 0) > mid:
                break
            for j in range(N):
                if f(i, j, 0) > mid:
                    break
                for k in range(N):
                    if f(i, j, k) <= mid:
                        count += 1
                        if count >= K:
                            return False
                    else:
                        break
        return count < K

    left, right = f(0, 0, 0), f(N - 1, N - 1, N - 1)
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            left = mid + 1
        else:
            right = mid - 1
    return left


if __name__ == "__main__":
    N, K = map(int, input().split())
    A = list(map(int, input().split()))
    B = list(map(int, input().split()))
    C = list(map(int, input().split()))

    print(kthLargestTriplet(A, B, C, K))
