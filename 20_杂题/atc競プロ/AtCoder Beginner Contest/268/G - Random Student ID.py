# 脑筋急转弯
# 字母表排列是随机的，有26！种情况，
# !每个人有一个由小写字母组成的编号，求在所有人编号中排名的期望
# !所有人的名字长度和<=5e5 (暗示前缀树插入次数O(S))

# https://zhuanlan.zhihu.com/p/563350153
# 你比对方大：
# （1）对方是你的一个前缀
# （2）你们没有共同的前缀，但是你第某个字母比对方第某个字母大
# !也就是说对方如果是我们的一个前缀，我们一定比他打，输的概率是1，不然呢，我们输的概率是1/2！
# !也就是说 我们只需要求一下是咱们前缀的数量乘1，
# !然后找出咱们不是对方前缀同时对方也不是咱们前缀的数乘1/2即可。
# 然后就是第一个rank是1而不是0。

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
INV2 = pow(2, MOD - 2, MOD)
from collections import defaultdict
from typing import Dict, Iterable, Optional


class TrieNodeWithParent:
    __slots__ = ("wordCount", "preCount", "children", "parent")

    def __init__(self, parent: Optional["TrieNodeWithParent"] = None):
        self.wordCount = 0  # 以当前节点为结尾的单词个数
        self.preCount = 0  # 以当前节点为前缀的单词个数
        self.parent = parent
        self.children: Dict[str, TrieNodeWithParent] = defaultdict(lambda: TrieNodeWithParent(self))


class TrieWithParent:
    __slots__ = "root"

    def __init__(self, iterable: Optional[Iterable[str]] = None):
        self.root = TrieNodeWithParent()
        for word in iterable or ():
            self.insert(word)

    def insert(self, word: str) -> None:
        if not word:
            return
        node = self.root
        for char in word:
            node = node.children[char]
            node.preCount += 1
        node.wordCount += 1

    def countWord(self, word: str) -> int:
        """树中有多少个单词word"""
        if not word:
            return 0
        node = self.root
        for char in word:
            if char not in node.children:
                return 0
            node = node.children[char]
        return node.wordCount

    def countWordStartsWith(self, prefix: str) -> int:
        """树中有多少个以prefix为前缀的单词"""
        if not prefix:
            return 0
        node = self.root
        for char in prefix:
            if char not in node.children:
                return 0
            node = node.children[char]
        return node.preCount

    def countWordAsPrefix(self, word: str) -> int:
        """树中有多少个单词是word的前缀"""
        node = self._find(word)
        if node is None:
            return 0

        res = 0
        while node.parent is not None:
            res += node.wordCount  # 重复的单词统计为多次而不是1次
            node = node.parent
        return res

    def _find(self, word: str) -> Optional[TrieNodeWithParent]:
        """返回word所在结点"""
        if not word:
            return None
        node = self.root
        for char in word:
            if char not in node.children:
                return None
            node = node.children[char]
        return node


n = int(input())
names = [input() for _ in range(n)]

trie = TrieWithParent(names)
for name in names:
    count1 = trie.countWordStartsWith(name)  # name是多少个单词的前缀
    count2 = trie.countWordAsPrefix(name)  # 有多少个单词是name的前缀
    other = n - 1 - count1 - count2  # 互相没有前缀关系
    print((1 + count2 + other * INV2) % MOD)
