# # !给定两个长度为n(n≤ 105）的多重集合(可以理解为Counter)A,B 。
# # 你可以执行以下操作多次
# # !1)选择A 中的一个数x ，把它变成2*x  (把一个数二进制后面添加0)
# # !2)选择A 中的一个数x ，把它变成[x/2] (抹去二进制下的末尾一位)
# # !现在问把A变成B最少需要多少次操作。如果不行输出-1

# # 我们可以从二进制的角度来理解，操作1是把一个数后面添加0，操作2是删除一个数的最后一个二进制数位。
# # 如果B 中的一个数可以通过A中的一个数通过操作1获得，那么我们也可以将B中的一个数，
# !通过删除末尾的0来获得A。

# 即:
# !A集合中，结点只有是左孩子才能向上跳
# !B集合中，结点都可以向上跳

# 如果A = 1001,Bi = 100100，我们可以通过2次操作使得B[i]变成A[i]。
# # !对于二进制问题我们可以使用01字典树


import sys
import os
from typing import List, Optional


sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


class XORTrieNode:
    __slots__ = "children", "counts", "char"

    def __init__(self, char: str):
        self.char = char  # 左右子树
        self.counts = [0, 0]
        self.children: List["XORTrieNode"] = [None, None]  # type: ignore


class XORTrie:
    def __init__(self):
        self.root = XORTrieNode("")
        self.res = 0

    def insert(self, bin: str, numsId: int) -> None:
        """
        注意这里插入与一般的01字典树不同,高位不补齐
        numsId表示当前数属于哪个数组

        注意0不插入
        """
        if bin == "0":
            return
        root = self.root
        for char in bin:
            bit = int(char)
            if root.children[bit] is None:
                root.children[bit] = XORTrieNode(char)
            root = root.children[bit]
        root.counts[numsId] += 1

    def dfs(self, cur: "XORTrieNode", parent: Optional["XORTrieNode"]) -> bool:
        """从叶子结点往上跳"""
        if not cur:
            return True
        for child in cur.children:
            if not self.dfs(child, cur):  # 先到底部
                return False

        # 消除相同的数
        min_ = min(cur.counts)
        cur.counts[0] -= min_
        cur.counts[1] -= min_

        # 不合法的情况 B结点有数但A中的结点不能向上跳(不为左子树结点)
        if cur.counts[1] and cur.char != "0":
            return False

        # 终点完全消除
        if parent is None:
            return cur.counts[0] == 0 and cur.counts[1] == 0

        # 统计要跳几次
        self.res += sum(cur.counts)
        parent.counts[0] += cur.counts[0]
        parent.counts[1] += cur.counts[1]
        return True


def main() -> None:
    n = int(input())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))

    trie = XORTrie()
    for num in nums1:
        trie.insert(bin(num)[2:] if num != 0 else "", 0)  # !注意0不插入
    for num in nums2:
        trie.insert(bin(num)[2:] if num != 0 else "", 1)

    if trie.dfs(trie.root, None):  # type: ignore
        print(trie.res)
    else:
        print(-1)


if __name__ == "__main__":

    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
