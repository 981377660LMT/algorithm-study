# オイラーのφ関数(Euler's-Phi-Function)
# 给定正整数n,[1,n]内的数与n互质的个数
# 欧拉函数phi函数:
#   φ(n) = n * (1 - 1/p1) * (1 - 1/p2) * ... * (1 - 1/pk)
#   其中p1,p2,...,pk为n的质因子


from typing import List


# O(nloglogn)
def getEulerPhiTable(n: int) -> List[int]:
    res = list(range(n + 1))
    for i in range(2, n + 1):
        if res[i] == i:
            for j in range(i, n + 1, i):
                res[j] = res[j] // i * (i - 1)
    return res


if __name__ == "__main__":
    n = int(input())
    print(getEulerPhiTable(n))
