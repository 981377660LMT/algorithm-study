# https://nyaannyaan.github.io/library/tree/pruefer-code.hpp
# https://www.luogu.com.cn/problem/solution/P6086

# Prufer序列-树上计数问题
# Prufer 序列和无根树的相互转化。

# Prufer序列可以将一个带标号n个节点的无根树用[1,n]中的n―2个整数表示，
# !即n个点的完全图的生成树与长度为n―2值域为[1, n]的数列构成的双射。
# Prufer序列可以方便的解决一类树相关的计数问题，比如凯莱定理:
# n个点的完全图的生成树有n**(n-2)个。

# 无根树转prufer算法流程O(n):
# !1.每次找到树中入度为 1 且编号最小的点。
# !2.删去该节点并在序列中添加与该节点相邻的结点编号。
# !3.重复上面的操作，知道剩下两个节点。

# prufer转无根树算法流程O(n):
# !1. 根据prufer序列的性质，可以算出每个点的度数。
# !2.每次选择一个度数为 1 的最小的结点编号。
# !3.将结点与当前枚举到的 Prüfer 序列的点连接，然后同时减掉两个点的度。
# !4.最后我们剩下两个度数为 1 的点，其中一个是结点 n。就把它们建立连接。

# 性质
# 1. 对应一颗有标号的无根树
# !2. 某编号结点的入度为在该点prufer序列中出现次数+1(没有出现的就是叶子结点)
# 3. 点数为n的有标号无根树个数为n**(n-2),有根树乘以n
# 4. 每个点度数为di的有标号无根树个数为(n-2)!/Π(di-1)!
# 5.在构造完 Prüfer 序列后原树中会剩下两个结点，其中一个一定是编号最大的点 n


from typing import List


# parents[i-1] 是以 n 为根时节点 i 的父节点 (1<=i<=n-1)
# n=6 parents=[3,6,4,6,1] 返回 [6, 1, 3, 4]
def parentsToPrufer(n: int, parents: List[int]) -> List[int]:
    parents = [-1] + parents
    deg = [0] * (n + 1)
    for i in range(1, n):
        deg[parents[i]] += 1
    prufer = [0] * (n - 2)
    i, j = 0, 1
    while i < n - 2:
        while deg[j] > 0:
            j += 1
        prufer[i] = parents[j]
        i += 1
        while i < n - 2:
            p = prufer[i - 1]
            deg[p] -= 1
            if p > j or deg[p] > 0:
                break
            prufer[i] = parents[p]
            i += 1
        j += 1
    return prufer


# n=6 prufer=[4,6,5,2] 返回 [4,6,6,5,2,0] 表示结点i的父节点为parents[i-1] (1<=i<=n-1)
def pruferToParents(n: int, prufer: List[int]) -> List[int]:
    deg = [0] * (n + 1)
    for p in prufer:
        deg[p] += 1
    prufer = prufer + [n]
    parents = [0] * (n + 1)
    i, j = 0, 1
    while i < n - 1:
        while deg[j] > 0:
            j += 1
        parents[j] = prufer[i]
        while i < n - 2:
            p = prufer[i]
            deg[p] -= 1
            if p > j or deg[p] > 0:
                break
            parents[p] = prufer[i + 1]
            i += 1
        i += 1
        j += 1
    return parents[1:]


# https://www.luogu.com.cn/problem/solution/P2290
# 给定每个点的度数，让你求有多少种符合条件的无根树。1<=n<=150
# (n-2)!/Π(di-1)!
def countTree(deg: "List[int]") -> int:
    n = len(deg)
    if n == 1:
        return 1 if deg[0] == 0 else 0
    fac = [1]
    for i in range(1, n + 5):
        fac.append(fac[-1] * i)
    res = fac[n - 2]
    degSum = 0
    for d in deg:
        if d == 0:  # 不连通
            return 0
        degSum += d - 1
        res //= fac[d - 1]
    return res if degSum == n - 2 else 0


if __name__ == "__main__":
    n = int(input())
    deg = list(map(int, input().split()))
    print(countTree(deg))
