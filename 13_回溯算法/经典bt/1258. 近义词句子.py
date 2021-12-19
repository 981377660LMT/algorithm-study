from typing import List
from collections import defaultdict


class UnionFind:
    def __init__(self):
        self.parent = defaultdict(str)

    def add(self, s: str) -> 'UnionFind':
        if s not in self.parent:
            self.parent[s] = s
        return self

    def find(self, x: str) -> str:
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x: str, y: str) -> None:
        rootX = self.find(x)
        rootY = self.find(y)
        self.parent[rootX] = rootY

    def isConnected(self, x: str, y: str) -> bool:
        return self.find(x) == self.find(y)


# 0 <= synonyms.length <= 10 => c^10
# synonyms[i].length == 2
class Solution:
    def generateSentences(self, synonyms: List[List[str]], text: str) -> List[str]:
        def bt(cur: int, path: List[str]):
            if cur == len(words):
                res.add(' '.join(path))
                return

            for s in adjMap[words[cur]] | {words[cur]}:
                path.append(s)
                bt(cur + 1, path)
                path.pop()

        # 处理近义词
        uf = UnionFind()
        for s1, s2 in synonyms:
            uf.add(s1).add(s2).union(s1, s2)

        # 构建邻接表
        # {"happy": ['happy', 'joy','glad','pleased']} 类似结构
        wordSet = set(sum(synonyms, []))
        adjMap = defaultdict(set)
        for w1 in wordSet:
            for w2 in wordSet:
                if uf.isConnected(w1, w2):
                    adjMap[w1].add(w2)

        # 回溯
        res = set()
        words = text.split(' ')
        bt(0, [])
        return sorted(res)


print(
    Solution().generateSentences(
        synonyms=[["happy", "joy"], ["sad", "sorrow"], ["joy", "cheerful"]],
        text="I am happy today but was sad yesterday",
    )
)

# 输出：
# ["I am cheerful today but was sad yesterday",
# "I am cheerful today but was sorrow yesterday",
# "I am happy today but was sad yesterday",
# "I am happy today but was sorrow yesterday",
# "I am joy today but was sad yesterday",
# "I am joy today but was sorrow yesterday"]
