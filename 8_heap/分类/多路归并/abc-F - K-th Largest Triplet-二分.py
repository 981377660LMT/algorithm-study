# F - K-th Largest Triplet
# https://atcoder.jp/contests/abc391/tasks/abc391_f
# https://atcoder.jp/contests/abc391/editorial/12086
# !给定三个长为N的数组A,B,C,求第K大的三元组A[i]*B[j]+B[j]*C[k]+C[k]*A[i].
# !N<=2e5,1<=K<=min(N,2e5)
# 二分+剪枝

import sys

input = lambda: sys.stdin.readline().rstrip("\r\n")


def kthLargestTriplet(A: list[int], B: list[int], C: list[int], K: int) -> int:
    A, B, C = sorted(A, reverse=True), sorted(B, reverse=True), sorted(C, reverse=True)  # type: ignore
    N = len(A)

    def f(i: int, j: int, k: int) -> int:
        return A[i] * B[j] + B[j] * C[k] + C[k] * A[i]

    def check(mid: int) -> bool:
        """
        不小于mid的数是否<k个.
        时间复杂度O(K).
        """
        count = 0
        for i in range(N):
            if f(i, 0, 0) < mid:
                break
            for j in range(N):
                if f(i, j, 0) < mid:
                    break
                for k in range(N):
                    if f(i, j, k) < mid:
                        break
                    count += 1
                    if count >= K:
                        return False
        return count < K

    left, right = f(N - 1, N - 1, N - 1), f(0, 0, 0)
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            right = mid - 1
        else:
            left = mid + 1
    return right


if __name__ == "__main__":
    N, K = map(int, input().split())
    A = list(map(int, input().split()))
    B = list(map(int, input().split()))
    C = list(map(int, input().split()))

    print(kthLargestTriplet(A, B, C, K))
