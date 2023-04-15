# 求出 h[x] = ∑f[i]*g[j] (0<=i,j<n,max(i,j)=x)


from typing import List


def bruteForce(f: List[int], g: List[int]) -> List[int]:
    n = len(f)
    h = [0] * n
    for i in range(n):
        for j in range(n):
            h[max(i, j)] += f[i] * g[j]
    return h


def zetaAndMobius(f: List[int], g: List[int]) -> List[int]:
    # 一维zeta:前缀和
    n = len(f)
    for i in range(n - 1):
        f[i + 1] += f[i]
        g[i + 1] += g[i]

    print(f, g)
    res = [a * b for a, b in zip(f, g)]
    for i in range(n - 1, 0, -1):
        res[i] -= res[i - 1]
    return f


print(bruteForce([1, 2, 3], [4, 5, 6]))  # [4, 23, 63]
print(zetaAndMobius([1, 2, 3], [4, 5, 6]))  # [4, 23, 63]
