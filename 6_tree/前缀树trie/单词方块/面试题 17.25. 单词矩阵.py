from collections import defaultdict
from typing import List, Dict


class TrieNode:
    __slots__ = ('wordCount', 'preCount', 'children')

    def __init__(self):
        self.wordCount = 0
        self.children: Dict[str, TrieNode] = defaultdict(TrieNode)


class Trie:
    def __init__(self):
        self.root = TrieNode()

    def insert(self, word: str) -> None:
        if not word:
            return
        node = self.root
        for char in word:
            node = node.children[char]
        node.wordCount += 1


# 创建由字母组成的面积最大的矩形，其中每一行组成一个单词(自左向右)，每一列也组成一个单词(自上而下)。
# words.length <= 1000
# words[i].length <= 100
class Solution:
    def maxRectangle(self, words: List[str]) -> List[str]:
        # 1.首先根据单词长度对字典进行分组，因为你知道每一列的长度必须相同，每一行的长度也必须相同。
        # 2.你能找到一个特定长宽的单词矩阵吗？如果尝试了所有的选项会怎样？
        # 3.当矩形看起来无效时，可以使用trie提前终止吗？

        def bt(row: int, path: List[str], roots: List[TrieNode]) -> None:
            """固定矩阵列的长度 对每一行选一个单词 检查每一列的前缀是否都在trie中"""
            nonlocal maxArea, res
            colLen = len(roots)
            if row * colLen > maxArea and all(root.wordCount > 0 for root in roots):
                maxArea = row * colLen
                res = path[:]

            for nextRow in lenMap[colLen]:
                if all(char in roots[i].children for i, char in enumerate(nextRow)):
                    nextRoots = [roots[i].children[char] for i, char in enumerate(nextRow)]
                    path.append(nextRow)
                    bt(row + 1, path, nextRoots)
                    path.pop()

        lenMap = defaultdict(set)
        for w in words:
            if w:
                lenMap[len(w)].add(w)
        allLen = sorted(lenMap, reverse=True)

        trie = Trie()
        for w in words:
            trie.insert(w)

        res, maxArea = [], 0
        # !从边长倒序出发，剪枝
        for colLen in allLen:
            roots = [trie.root] * colLen
            bt(0, [], roots)
            if maxArea >= colLen * colLen:
                break

        return res
