# F - K-th Largest Triplet
# https://atcoder.jp/contests/abc391/tasks/abc391_f
# !给定三个长为N的数组A,B,C,求第K大的三元组A[i]*B[j]+B[j]*C[k]+C[k]*A[i].
# N<=2e5,1<=K<=min(N,2e5)
# 多路归并+堆
# https://atcoder.jp/contests/abc391/editorial/12105 多路归并防止重入堆的方法

from heapq import heappop, heappush
import sys

input = lambda: sys.stdin.readline().rstrip("\r\n")


def kthLargestTriplet(A: list[int], B: list[int], C: list[int], K: int) -> int:
    A, B, C = sorted(A, reverse=True), sorted(B, reverse=True), sorted(C, reverse=True)  # type: ignore

    N = len(A)
    visited = set()
    pq = []  # 大根堆

    def add(i: int, j: int, k: int) -> None:
        if i >= N or j >= N or k >= N:
            return
        h = i * N * N + j * N + k
        if h in visited:
            return
        visited.add(h)
        v = A[i] * B[j] + B[j] * C[k] + C[k] * A[i]
        heappush(pq, (-v, i, j, k))

    add(0, 0, 0)
    for t in range(K):
        v, i, j, k = heappop(pq)
        if t == K - 1:
            return -v
        add(i + 1, j, k)
        add(i, j + 1, k)
        add(i, j, k + 1)

    raise ValueError("unreachable")


if __name__ == "__main__":
    N, K = map(int, input().split())
    A = list(map(int, input().split()))
    B = list(map(int, input().split()))
    C = list(map(int, input().split()))

    print(kthLargestTriplet(A, B, C, K))
