# オイラーのφ関数(Euler's-Phi-Function)
# 给定正整数n,[1,n]内的数与n互质的个数
# 欧拉函数phi函数:
#   φ(n) = n * (1 - 1/p1) * (1 - 1/p2) * ... * (1 - 1/pk)
#   其中p1,p2,...,pk为n的质因子


# O(sqrt(n))
def eulerPhi(n: int) -> int:
    res = n
    upper = int(n**0.5)
    for i in range(2, upper + 1):
        if n % i == 0:
            res -= res // i
            while n % i == 0:
                n //= i
    if n > 1:
        res -= res // n
    return res


if __name__ == "__main__":
    n = int(input())
    print(eulerPhi(n))
