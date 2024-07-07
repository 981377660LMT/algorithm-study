# 100350. 最小代价构造字符串

# 给你一个字符串 target、一个字符串数组 words 以及一个整数数组 costs，这两个数组长度相同。

# 设想一个空字符串 s。

# 你可以执行以下操作任意次数（包括零次）：

# 选择一个在范围  [0, words.length - 1] 的索引 i。
# 将 words[i] 追加到 s。
# 该操作的成本是 costs[i]。
# 返回使 s 等于 target 的 最小 成本。如果不可能，返回 -1。

from typing import Generator, Generic, Iterable, List, Tuple, TypeVar

INF = int(2e18)


T = TypeVar("T", str, int)


class ACAutoMatonMap(Generic[T]):
    """
    不调用 BuildSuffixLink 就是Trie, 调用 BuildSuffixLink 就是AC自动机.
    每个状态对应Trie中的一个结点, 也对应一个字符串.
    """

    __slots__ = ("wordPos", "children", "_link", "_linkWord", "_bfsOrder")

    def __init__(self):
        self.wordPos = []
        """wordPos[i] 表示加入的第i个模式串对应的节点编号."""
        self.children = [{}]
        """children[v][c] 表示节点v通过字符c转移到的节点."""
        self._link = []
        """又叫fail.指向当前节点最长真后缀对应结点,例如"bc"是"abc"的最长真后缀."""
        self._linkWord = []
        self._bfsOrder = []
        """结点的拓扑序,0表示虚拟节点."""

    def addString(self, string: Iterable[T]) -> int:
        if not string:
            return 0
        pos = 0
        for char in string:
            nexts = self.children[pos]
            if char in nexts:
                pos = nexts[char]
            else:
                nextState = len(self.children)
                nexts[char] = nextState
                pos = nextState
                self.children.append({})
        self.wordPos.append(pos)
        return pos

    def addChar(self, pos: int, char: T) -> int:
        nexts = self.children[pos]
        if char in nexts:
            return nexts[char]
        nextState = len(self.children)
        nexts[char] = nextState
        self.children.append({})
        return nextState

    def move(self, pos: int, char: T) -> int:
        children, link = self.children, self._link
        while True:
            nexts = children[pos]
            if char in nexts:
                return nexts[char]
            if pos == 0:
                return 0
            pos = link[pos]

    def buildSuffixLink(self):
        """
        构建后缀链接(失配指针).
        """
        self._link = [-1] * len(self.children)
        self._bfsOrder = [0] * len(self.children)
        head, tail = 0, 1
        while head < tail:
            v = self._bfsOrder[head]
            head += 1
            for char, next_ in self.children[v].items():
                self._bfsOrder[tail] = next_
                tail += 1
                f = self._link[v]
                while f != -1 and char not in self.children[f]:
                    f = self._link[f]
                self._link[next_] = f
                if f == -1:
                    self._link[next_] = 0
                else:
                    self._link[next_] = self.children[f][char]

    def linkWord(self, pos: int) -> int:
        """
        `linkWord`指向当前节点的最长后缀对应的节点.
        区别于`_link`,`linkWord`指向的节点对应的单词不会重复.
        即不会出现`_link`指向某个长串局部的恶化情况.

        时间复杂度 O(sqrt(n)).
        """
        if len(self._linkWord) == 0:
            hasWord = [False] * len(self.children)
            for v in self.wordPos:
                hasWord[v] = True
            self._linkWord = [0] * len(self.children)
            link, linkWord = self._link, self._linkWord
            for v in self._bfsOrder:
                if v != 0:
                    p = link[v]
                    linkWord[v] = p if hasWord[p] else linkWord[p]
        return self._linkWord[pos]

    def getCounter(self) -> List[int]:
        """
        获取每个状态包含的模式串的个数.
        时空复杂度 O(n).
        """
        counter = [0] * len(self.children)
        for pos in self.wordPos:
            counter[pos] += 1
        for v in self._bfsOrder:
            if v != 0:
                counter[v] += counter[self._link[v]]
        return counter

    def getIndexes(self) -> List[List[int]]:
        """
        获取每个状态包含的模式串的索引(有序).
        时空复杂度 O(nsqrtn).
        """
        res = [[] for _ in range(len(self.children))]
        for i, pos in enumerate(self.wordPos):
            res[pos].append(i)
        for v in self._bfsOrder:
            if v != 0:
                from_, _children = self._link[v], v
                arr1, arr2 = res[from_], res[_children]
                arr3 = []
                i, j = 0, 0
                while i < len(arr1) and j < len(arr2):
                    if arr1[i] < arr2[j]:
                        arr3.append(arr1[i])
                        i += 1
                    elif arr1[i] > arr2[j]:
                        arr3.append(arr2[j])
                        j += 1
                    else:
                        arr3.append(arr1[i])
                        i += 1
                        j += 1
                arr3 += arr1[i:] + arr2[j:]
                res[_children] = arr3
        return res

    def dp(self) -> Generator[Tuple[int, int], None, None]:
        for v in self._bfsOrder:
            if v != 0:
                yield self._link[v], v

    def buildFailTree(self) -> List[List[int]]:
        adjList = [[] for _ in range(len(self.children))]
        for v in self._bfsOrder:
            if v != 0:
                adjList[self._link[v]].append(v)
        return adjList

    def buildTrieTree(self) -> List[List[int]]:
        adjList = [[] for _ in range(len(self.children))]

        def dfs(pos: int) -> None:
            for next_ in self.children[pos].values():
                adjList[pos].append(next_)
                dfs(next_)

        dfs(0)
        return adjList

    def search(self, string: Iterable[T]) -> int:
        """返回string在trie树上的节点位置.如果不存在,返回0."""
        if not string:
            return 0
        pos = 0
        for char in string:
            if pos < 0 or pos >= len(self.children):
                return 0
            nexts = self.children[pos]
            if char in nexts:
                pos = nexts[char]
            else:
                return 0
        return pos

    def empty(self) -> bool:
        return len(self.children) == 1

    def clear(self) -> None:
        self.wordPos = []
        self.children = [{}]
        self._link = []
        self._linkWord = []
        self._bfsOrder = []

    @property
    def size(self) -> int:
        return len(self.children)

    def __len__(self) -> int:
        return len(self.children)


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def minimumCost(self, target: str, words: List[str], costs: List[int]) -> int:
        acm = ACAutoMatonMap()
        for word in words:
            acm.addString(word)
        acm.buildSuffixLink()
        dp = [INF] * (len(target) + 1)
        dp[0] = 0
        pos = 0
        indexes = acm.getIndexes()
        for i, c in enumerate(target):
            pos = acm.move(pos, c)
            for wordIndex in indexes[pos]:
                wordLen = len(words[wordIndex])
                if i + 1 >= wordLen:
                    dp[i + 1] = min2(dp[i + 1], dp[i + 1 - wordLen] + costs[wordIndex])

        return dp[len(target)] if dp[len(target)] != INF else -1


# "abcdef"
# ["abdef","abc","d","def","ef"]
# [100,1,1,10,5]

print(Solution().minimumCost("abcdef", ["abdef", "abc", "d", "def", "ef"], [100, 1, 1, 10, 5]))
