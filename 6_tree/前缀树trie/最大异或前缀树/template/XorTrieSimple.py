from typing import List, Tuple


class XorTrieSimpleNode:
    __slots__ = "count", "zero", "one"

    def __init__(self):
        self.count = 0
        self.zero = None
        self.one = None


class XorTrieSimple:
    __slots__ = "root", "bit"

    def __init__(self, upper: int):
        self.root = XorTrieSimpleNode()
        self.bit = upper.bit_length()

    def insert(self, num: int) -> "XorTrieSimpleNode":
        root = self.root
        for i in range(self.bit, -1, -1):
            bit = (num >> i) & 1
            if bit:
                if root.one is None:
                    root.one = XorTrieSimpleNode()
                root = root.one
            else:
                if root.zero is None:
                    root.zero = XorTrieSimpleNode()
                root = root.zero
            root.count += 1
        return root

    def remove(self, num: int) -> None:
        """需要保证num在trie中."""
        root = self.root
        for i in range(self.bit, -1, -1):
            bit = (num >> i) & 1
            root = bit.one if bit else root.zero  # type: ignore
            root.count -= 1  # type: ignore

    def query(self, num: int) -> int:
        """查询num与trie中的最大异或值."""
        root = self.root
        res = 0
        for i in range(self.bit, -1, -1):
            if root is None:  # type: ignore
                break
            bit = (num >> i) & 1
            if bit:
                if root.zero is not None and root.zero.count > 0:
                    res |= 1 << i
                    root = root.zero
                else:
                    root = root.one
            else:
                if root.one is not None and root.one.count > 0:
                    res |= 1 << i
                    root = root.one
                else:
                    root = root.zero

        return res


if __name__ == "__main__":
    # https://leetcode.cn/problems/maximum-strong-pair-xor-ii/
    # 给你一个下标从 0 开始的整数数组 nums 。如果一对整数 x 和 y 满足以下条件，则称其为 强数对 ：
    # |x - y| <= min(x, y)
    # 你需要从 nums 中选出两个整数，且满足：这两个整数可以形成一个强数对，并且它们的按位异或（XOR）值是在该数组所有强数对中的 最大值 。
    # 返回数组 nums 所有可能的强数对中的 最大 异或值。
    #
    # 排序后，变为 y <= 2x
    class Solution:
        def maximumStrongPairXor(self, nums: List[int]) -> int:
            nums.sort()
            res, left, n = 0, 0, len(nums)
            trie = XorTrieSimple(max(nums))
            for right in range(n):
                cur = nums[right]
                trie.insert(cur)
                while left <= right and cur > 2 * nums[left]:
                    trie.remove(nums[left])
                    left += 1
                res = max(res, trie.query(cur))
            return res

    # https://leetcode.cn/problems/maximum-xor-of-two-numbers-in-an-array/description/
    class Solution2:
        def findMaximumXOR(self, nums: List[int]) -> int:
            trie = XorTrieSimple(max(nums))
            res = 0
            for num in nums:
                trie.insert(num)
                res = max(res, trie.query(num))
            return res

    # P4551 最长异或路径
    # 给定一棵 n 个点的带权树，结点下标从 0 开始到 n。
    # 寻找树中找两个结点，求最长的异或路径。
    # 异或路径指的是指两个结点之间唯一路径上的所有边权的异或。
    #
    # 树上差分
    # 树上 x 到 y 的路径上所有边权的 xor 结果就等于 `D[x] xor D[y]`。
    # !其中 D[x]表示根节点到 x 的异或值,重叠路径抵消了(前缀异或)
    # 所以，`问题就变成了从 D[1]~D[N]这 N 个数中选出两个，xor 的结果最大`
    # 时间复杂度O(n)

    # https://www.luogu.com.cn/problem/P4551
    def solve(n: int, edges: List[Tuple[int, int, int]]) -> int:
        adjList = [[] for _ in range(n)]
        for u, v, w in edges:
            adjList[u].append((v, w))
            adjList[v].append((u, w))

        xorToRoot = []

        def dfs(cur: int, pre: int, preXor: int) -> int:
            xorToRoot.append(preXor)
            for v, w in adjList[cur]:
                if v != pre:
                    dfs(v, cur, preXor ^ w)

        dfs(0, -1, 0)

        trie = XorTrieSimple(max(xorToRoot, default=1))
        res = 0
        for xor in xorToRoot:
            trie.insert(xor)
            res = max(res, trie.query(xor))
        return res

    n = int(input())
    edges = []
    for _ in range(n - 1):
        u, v, w = map(int, input().split())
        edges.append((u - 1, v - 1, w))
    print(solve(n, edges))
