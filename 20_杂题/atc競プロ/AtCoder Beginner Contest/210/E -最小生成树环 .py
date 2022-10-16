# n个点，m种边，每种边都有一个offset,weight ,
# 即你可以选择任意一个点x和(x+offset)%n连一条长度为weight的边,
# !每种边都可以使用任意次，顺序也可以任意，问把n个点连成一个连通图的最小权值是多少
# 如果不能连成一个连通图，输出-1
# n<=1e9 m<=1e5 (点数巨大)

# n太大 显然不能一条边一条边地合并
# !需要借助数学加速合并
# !每个回合,一种边能连接的边的数量是多少?(或者说,每个回合有多少个连通分量?/每次减少多少个连通分量?)
# 如果u和v可以通过前t条边相连，也就是存在一组正整数解(k1,k2,...,kt)满足
# !v ≡ u + k1 * offset1 + k2 * offset2 + ... + kt * offsett (mod n)
# 即 k0*n+k1*offset1+k2*offset2+...+kt*offsett=v-u 有整数解(k0,k1,k2,...,kt)
# !根据裴蜀定理 gcd(n,offset1,offset2,...,offsett) | v-u
# 因此有gcd(n,offset1,offset2,...,offsett)组连通分量


# !因此每回合合并前联通分量个数为preGcd , 当前回合边合并后联通分量个数为 `gcd(preGcd,offseti)`

from math import gcd
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, m = map(int, input().split())
    edges = [tuple(map(int, input().split())) for _ in range(m)]  # offset,weight
    edges.sort(key=lambda x: x[1])  # 按权值从小到大排序

    res, preGcd = 0, n
    for offset, weight in edges:
        curGcd = gcd(preGcd, offset)
        res += weight * (preGcd - curGcd)
        preGcd = curGcd

    print(res if preGcd == 1 else -1)
