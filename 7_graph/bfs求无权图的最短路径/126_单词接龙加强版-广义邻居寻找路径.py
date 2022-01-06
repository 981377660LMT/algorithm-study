from collections import defaultdict, deque
from typing import List

INF = 0x7FFFFFFF


class Solution:
    def findLadders(self, beginWord: str, endWord: str, wordList: List[str]) -> List[List[str]]:
        dist = defaultdict(lambda: INF)
        dist[beginWord] = 0
        adjMap = defaultdict(list)
        res = []

        for word in wordList:
            for i in range(len(word)):
                adjMap[word[:i] + "*" + word[i + 1 :]].append(word)
        queue = deque([[beginWord]])

        while queue:
            path = queue.popleft()
            cur = path[-1]
            if cur == endWord:
                res.append(path.copy())
            else:
                for i in range(len(cur)):
                    for next in adjMap[cur[:i] + "*" + cur[i + 1 :]]:
                        if dist[cur] + 1 <= dist[next]:
                            queue.append(path + [next])
                            dist[next] = dist[cur] + 1
        return res


print(
    Solution().findLadders(
        beginWord="hit", endWord="cog", wordList=["hot", "dot", "dog", "lot", "log", "cog"]
    )
)

