# !ACAutoMatonMap 更快.

from typing import Generator, List, Tuple

INF = int(1e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


class ACAutoMatonArray:
    """
    不调用`BuildSuffixLink`就是Trie,调用`BuildSuffixLink`就是AC自动机.
    每个状态对应Trie中的一个结点,也对应一个前缀.
    """

    __slots__ = (
        "wordPos",
        "parent",
        "children",
        "_bfsOrder",
        "_link",
        "_linkWord",
        "_sigma",
        "_offset",
        "_needUpdateChildren",
    )

    def __init__(self, sigma: int = 26, offset: int = 97):
        self.wordPos = []
        self.parent = []
        self.children = []
        self._bfsOrder = []
        self._link = []
        self._linkWord = []
        self._sigma = sigma
        self._offset = offset
        self._needUpdateChildren = False
        self.clear()

    def addString(self, string: str) -> int:
        """添加一个字符串，返回最后一个字符对应的节点编号."""
        if not string:
            return 0
        pos = 0
        for c in string:
            ord_ = ord(c) - self._offset
            if self.children[pos][ord_] == -1:
                self.children[pos][ord_] = self._newNode()
                self.parent[-1] = pos
            pos = self.children[pos][ord_]
        self.wordPos.append(pos)
        return pos

    def addChar(self, pos: int, char: str) -> int:
        """在pos位置添加一个字符,返回新的节点编号."""
        ord_ = ord(char) - self._offset
        if self.children[pos][ord_] != -1:
            return self.children[pos][ord_]
        self.children[pos][ord_] = self._newNode()
        self.parent[-1] = pos
        return self.children[pos][ord_]

    def buildSuffixLink(self, needUpdateChildren: bool = True):
        """
        构建后缀链接(失配指针).
        needUpdateChildren 表示是否需要更新children数组(连接trie图).
        """
        self._needUpdateChildren = needUpdateChildren
        self._link = [-1] * len(self.children)
        self._bfsOrder = [0] * len(self.children)
        link, order, children = self._link, self._bfsOrder, self.children
        head, tail = 0, 0
        order[tail] = 0
        tail += 1
        while head < tail:
            v = order[head]
            head += 1
            for i, next_ in enumerate(children[v]):
                if next_ == -1:
                    continue
                order[tail] = next_
                tail += 1
                f = link[v]
                while f != -1 and children[f][i] == -1:
                    f = link[f]
                if f == -1:
                    link[next_] = 0
                else:
                    link[next_] = children[f][i]
        if not needUpdateChildren:
            return
        for v in order:
            for i, next_ in enumerate(children[v]):
                if next_ == -1:
                    f = link[v]
                    if f == -1:
                        children[v][i] = 0
                    else:
                        children[v][i] = children[f][i]

    def move(self, pos: int, char: str) -> int:
        """pos: DFA的状态集, char: DFA的字符集."""
        ord_ = ord(char) - self._offset
        if self._needUpdateChildren:
            return self.children[pos][ord_]
        while True:
            nexts = self.children[pos]
            if nexts[ord_] != -1:
                return nexts[ord_]
            if pos == 0:
                return 0
            pos = self._link[pos]

    def linkWord(self, pos: int) -> int:
        """
        `linkWord`指向当前节点的最长后缀对应的节点.
        区别于`_link`,`linkWord`指向的节点对应的单词不会重复.
        即不会出现`_link`指向某个长串局部的恶化情况.
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

    def clear(self) -> None:
        self.wordPos.clear()
        self.parent.clear()
        self.children.clear()
        self._link.clear()
        self._linkWord.clear()
        self._bfsOrder.clear()
        self._newNode()

    def getCounter(self) -> List[int]:
        """
        获取每个状态包含的模式串的个数.
        时空复杂度 O(n).
        """
        counter = [0] * len(self.children)
        for pos in self.wordPos:
            counter[pos] += 1
        link = self._link
        for v in self._bfsOrder:
            if v != 0:
                counter[v] += counter[link[v]]
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
        """按照拓扑序进行转移, 每次返回(link, cur)."""
        link = self._link
        for v in self._bfsOrder:
            if v != 0:
                yield link[v], v

    def search(self, string: str) -> int:
        """返回string在trie树上的节点位置.如果不存在,返回0."""
        if not string:
            return 0
        pos = 0
        n = self.size
        for c in string:
            if pos >= n or pos < 0:
                return 0
            ord_ = ord(c) - self._offset
            next_ = self.children[pos][ord_]
            if next_ == -1:
                return 0
            pos = next_
        return pos

    def buildFailTree(self) -> List[List[int]]:
        res = [[] for _ in range(len(self.children))]
        link = self._link
        for v in self._bfsOrder:
            if v != 0:
                res[link[v]].append(v)
        return res

    def buildTrieTree(self) -> List[List[int]]:
        res = [[] for _ in range(len(self.children))]
        for i, v in enumerate(self.parent[1:], 1):
            res[v].append(i)
        return res

    def empty(self) -> bool:
        return len(self.children) == 1

    @property
    def size(self) -> int:
        """自动机中的节点(状态)数量,包括虚拟节点0."""
        return len(self.children)

    def __len__(self) -> int:
        return len(self.children)

    def _newNode(self) -> int:
        self.parent.append(-1)
        self.children.append([-1] * self._sigma)
        return len(self.children) - 1


if __name__ == "__main__":

    class Solution:
        # 100350. 最小代价构造字符串
        # https://leetcode.cn/problems/construct-string-with-minimum-cost/description/
        def minimumCost(self, target: str, words: List[str], costs: List[int]) -> int:
            acm = ACAutoMatonArray(sigma=26, offset=97)
            for word in words:
                acm.addString(word)
            acm.buildSuffixLink(True)

            nodeCosts, nodeDepth = [INF] * acm.size, [0] * acm.size
            for i, pos in enumerate(acm.wordPos):
                nodeCosts[pos] = min2(nodeCosts[pos], costs[i])
                nodeDepth[pos] = len(words[i])

            n = len(target)
            dp = [INF] * (n + 1)
            dp[0] = 0
            pos = 0
            for i, char in enumerate(target):
                pos = acm.move(pos, char)
                cur = pos
                while cur:
                    dp[i + 1] = min2(dp[i + 1], dp[i + 1 - nodeDepth[cur]] + nodeCosts[cur])
                    cur = acm.linkWord(cur)
            return dp[n] if dp[n] != INF else -1

        # 2781. 最长合法子字符串的长度
        # https://leetcode.cn/problems/length-of-the-longest-valid-substring/
        def longestValidSubstring(self, word: str, forbidden: List[str]) -> int:
            acm = ACAutoMatonArray(sigma=26, offset=97)
            for s in forbidden:
                acm.addString(s)
            acm.buildSuffixLink(True)

            minWordLen = [INF] * acm.size  # 每个状态匹配到的模式串的最小长度
            for i, pos in enumerate(acm.wordPos):
                minWordLen[pos] = len(forbidden[i])
            for link, cur in acm.dp():
                minWordLen[cur] = min2(minWordLen[cur], minWordLen[link])

            res, left, pos = 0, 0, 0
            for right, char in enumerate(word):
                pos = acm.move(pos, char)
                left = max2(left, right - minWordLen[pos] + 2)
                res = max2(res, right - left + 1)
            return res

        # 面试题 17.17. 多次搜索
        # https://leetcode.cn/problems/multi-search-lcci/
        def multiSearch(self, big: str, smalls: List[str]) -> List[List[int]]:
            """多模式匹配indexOfAll."""
            acm = ACAutoMatonArray()
            for s in smalls:
                acm.addString(s)
            acm.buildSuffixLink()

            matchIndexes = acm.getIndexes()
            res = [[] for _ in range(len(smalls))]
            pos = 0
            for i, char in enumerate(big):
                pos = acm.move(pos, char)
                for index in matchIndexes[pos]:
                    res[index].append(i - len(smalls[index]) + 1)
            return res

        # 1032. 字符流
        # https://leetcode.cn/problems/stream-of-characters/
        class StreamChecker:
            def __init__(self, words: List[str]):
                self._acm = ACAutoMatonArray()
                for w in words:
                    self._acm.addString(w)
                self._acm.buildSuffixLink()
                self._counter = self._acm.getCounter()
                self._pos = 0

            def query(self, letter: str) -> bool:
                """从字符流中接收一个新字符，如果字符流中的任一非空后缀能匹配 words 中的某一字符串，返回 true"""
                self._pos = self._acm.move(self._pos, letter)
                return self._counter[self._pos] > 0

    # "twnpxyhva"
    # ["pxyhva","twnpxyhva","wnpx"]
    # [3,19,8]
    print(Solution().minimumCost("twnpxyhva", ["pxyhva", "twnpxyhva", "wnpx"], [3, 19, 8]) == 8)
