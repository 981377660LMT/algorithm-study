# 最大异或路径 - 异或前缀和
# 144. 最长异或值路径
# 给定一个树，树上的边都具有权值。
# 树中一条路径的异或长度被定义为路径上所有边的权值的异或和：
# 给定上述的具有 n 个节点的树，你能找到异或长度最大的路径吗？
# 异或最大的路径

# 树上 x 到 y 的路径上所有边权的 xor 结果就等于 `D[x] xor D[y]`。
# 其中 D[x]表示根节点到 x 的异或值,重叠路径抵消了(前缀异或)
# 所以，`问题就变成了从 D[1]~D[N]这 N 个数中选出两个，xor 的结果最大`
# 时间复杂度O(n)

# https://www.acwing.com/problem/content/description/146/
from typing import List
from collections import defaultdict, namedtuple


def useXorTrie(bitLength=31):
    trieRoot = [None, None, 0]

    def insert(num: int) -> None:
        root = trieRoot
        for i in range(bitLength, -1, -1):
            bit = (num >> i) & 1
            if root[bit] is None:
                root[bit] = [None, None, 0]
            root[bit][2] += 1
            root = root[bit]

    def search(num: int) -> int:
        root = trieRoot
        res = 0
        for i in range(bitLength, -1, -1):
            if root is None:  # Trie中未插入
                break

            bit = (num >> i) & 1
            needBit = bit ^ 1
            if root[needBit] is not None and root[needBit][2] > 0:
                res = res << 1 | 1
                root = root[needBit]
            else:
                res = res << 1
                root = root[bit]

        return res

    def discard(num: int) -> None:
        root = trieRoot
        for i in range(bitLength, -1, -1):
            if root is None:  # Trie中未插入
                break

            bit = (num >> i) & 1
            if root[bit] is not None:
                root[bit][2] -= 1
            root = root[bit]

    return namedtuple('XorTrie', ['insert', 'search', 'discard'])(insert, search, discard)


def maxXor(nums: List[int]) -> int:
    res = 0
    xorTrie = useXorTrie()
    for num in nums:
        res = max(res, xorTrie.search(num))
        xorTrie.insert(num)
    return res


def dfs(cur: int, parent: int, curXor: int) -> None:
    for next, weight in adjMap[cur].items():
        if next == parent:
            continue
        xors[next] = curXor ^ weight
        dfs(next, cur, xors[next])


adjMap = defaultdict(lambda: defaultdict(int))
n = int(input())
for _ in range(n - 1):
    u, v, w = map(int, input().split())
    adjMap[u][v] = w
    adjMap[v][u] = w

xors = [0] * n
dfs(0, -1, 0)
print(maxXor(xors))
