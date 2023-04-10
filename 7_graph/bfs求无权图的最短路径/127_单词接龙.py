# !找出并返回所有从 beginWord 到 endWord 的 最短距离
# wordList 中的所有单词 互不相同


from collections import deque
from typing import List


class Solution:
    def ladderLength(self, beginWord: str, endWord: str, wordList: List[str]) -> int:
        def genNexts(cur: str):
            for i in range(len(cur)):
                for j in range(26):
                    next = cur[:i] + chr(97 + j) + cur[i + 1 :]  # !字符串哈希优化
                    if next in S and next != cur:
                        yield next

        S = set(wordList)
        if endWord not in S:  # 不存在转换序列
            return 0
        queue = deque([(beginWord, 1)])
        visited = set([beginWord])
        while queue:
            cur, dist = queue.popleft()
            if cur == endWord:
                return dist
            for next in genNexts(cur):
                if next not in visited:
                    queue.append((next, dist + 1))
                    visited.add(next)
        return 0  # 不存在
