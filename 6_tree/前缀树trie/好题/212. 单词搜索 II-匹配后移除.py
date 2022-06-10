# 给定一个 m x n 二维字符网格 board 和一个单词（字符串）列表 words， 返回所有二维网格上的单词 。
# 1 <= m, n <= 12
# 1 <= words.length <= 3 * 1e4
# 1 <= words[i].length <= 10


from collections import defaultdict
from typing import Dict, List


class TrieNode:
    __slots__ = ('wordCount', 'preCount', 'children', 'word')

    def __init__(self):
        self.wordCount = 0
        self.preCount = 0
        self.word = ''  # 存储当前的单词
        self.children: Dict[str, TrieNode] = defaultdict(TrieNode)


class Trie:
    def __init__(self):
        self.root = TrieNode()  # 声明为公有 便于改造

    def insert(self, word: str) -> None:
        if not word:
            return
        node = self.root
        for char in word:
            if char not in node.children:
                node.children[char] = TrieNode()
            node = node.children[char]
            node.preCount += 1
        node.wordCount += 1
        node.word = word

    def countWord(self, word: str) -> int:
        """树中是否存在word,返回个数"""
        if not word:
            return 0
        node = self.root
        for char in word:
            if char not in node.children:
                return 0
            node = node.children[char]
        return node.wordCount

    def countPre(self, prefix: str) -> int:
        """树中是否存在以prefix为前缀的单词,返回个数"""
        if not prefix:
            return 0
        node = self.root
        for char in prefix:
            if char not in node.children:
                return 0
            node = node.children[char]
        return node.preCount

    def discard(self, word: str) -> None:
        """从树中删除1个word 必须保证树中存在word"""
        if not word:
            return
        node = self.root
        for char in word:
            if char not in node.children:
                return
            node = node.children[char]
            node.preCount -= 1
        node.wordCount -= 1


class Solution:
    def findWords(self, board: List[List[str]], words: List[str]) -> List[str]:
        """
        1. 我们可以将匹配到的单词从前缀树中移除，来避免重复寻找相同的单词
        2. 弹出 '#' 后，非 '#' 的叶子结点也可以弹出
        """

        def bt(root: TrieNode, row: int, col: int) -> None:
            if board[row][col] not in root.children:
                return

            cur = board[row][col]
            board[row][col] = '#'

            nextNode = root.children[cur]
            if nextNode.wordCount > 0:
                res.add(nextNode.word)
                nextNode.wordCount = -1  # !直接移除 无需去重了

            for nr, nc in [(row - 1, col), (row + 1, col), (row, col - 1), (row, col + 1)]:
                if 0 <= nr < ROW and 0 <= nc < COL:
                    bt(nextNode, nr, nc)

            board[row][col] = cur

        trie = Trie()
        for w in words:
            trie.insert(w)

        res = set()
        ROW, COL = len(board), len(board[0])
        for row in range(ROW):
            for col in range(COL):
                bt(trie.root, row, col)

        return list(res)

