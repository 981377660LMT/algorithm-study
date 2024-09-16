# abc366-G - XOR Neighbors
# https://atcoder.jp/contests/abc366/tasks/abc366_g
# 给图上的每个点赋一个1~2^60-1的点权，使得对所有节点的邻居xor和为0，或报告无解。
# n<=60


import sys

readline = sys.stdin.readline


def get_basis_and_permutation(A):
    n = len(A)
    P = []
    basis = []
    zero = []
    for i in range(n):
        v = A[i]
        for b, j in basis:
            if v > v ^ b:
                v ^= b
                P.append((i, j))
        if v:
            basis.append((v, i))
        else:
            zero.append(i)
    return basis, P, zero


def solve_linear_equation_F2_row(A, Y):
    basis, P, zero = get_basis_and_permutation(A)
    X = [0] * len(A)
    for b, i in basis:
        top = b.bit_length() - 1
        if Y >> top & 1:
            X[i] = 1
            Y ^= b
    if Y:
        return [-1], []
    for i, j in P[::-1]:
        X[j] ^= X[i]

    kernel = []
    for zi in zero:
        ker = [0] * len(A)
        ker[zi] = 1
        for i, j in P[::-1]:
            ker[j] ^= ker[i]
        kernel.append(ker)
    return X, kernel


n, m = map(int, readline().split())
uv = [list(map(int, readline().split())) for _ in range(m)]


A = [0] * n
for u, v in uv:
    u -= 1
    v -= 1
    A[u] |= 1 << v
    A[v] |= 1 << u


X, ker = solve_linear_equation_F2_row(A, 0)

res = [0] * n
for i, lst in enumerate(ker):
    for j in range(n):
        res[j] |= lst[j] << i

if any(i == 0 for i in res):
    print("No")
else:
    print("Yes")
    print(*res)
