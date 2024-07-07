from typing import Generator, List, Tuple

INF = int(1e18)


class ACAutoMatonArray:
    """
    不调用`BuildSuffixLink`就是Trie,调用`BuildSuffixLink`就是AC自动机.
    每个状态对应Trie中的一个结点,也对应一个前缀.
    """

    __slots__ = (
        "wordPos",
        "parent",
        "link",
        "children",
        "_bfsOrder",
        "_sigma",
        "_offset",
        "_needUpdateChildren",
    )

    def __init__(self, sigma: int = 26, offset: int = 97):
        self.wordPos = []
        self.parent = []
        self.link = []
        self.children = []
        self._bfsOrder = []
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

    def addChar(self, pos: int, ord: int) -> int:
        """在pos位置添加一个字符,返回新的节点编号."""
        ord -= self._offset
        if self.children[pos][ord] != -1:
            return self.children[pos][ord]
        self.children[pos][ord] = self._newNode()
        self.parent[-1] = pos
        return self.children[pos][ord]

    def buildSuffixLink(self, needUpdateChildren: bool = True):
        """
        构建后缀链接(失配指针).
        needUpdateChildren 表示是否需要更新children数组(连接trie图).
        """
        self._needUpdateChildren = needUpdateChildren
        self.link = [-1] * len(self.children)
        self._bfsOrder = [0] * len(self.children)
        link, order, children = self.link, self._bfsOrder, self.children
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

    def move(self, pos: int, ord: int) -> int:
        """pos: DFA的状态集, ord: DFA的字符集."""
        ord -= self._offset
        if self._needUpdateChildren:
            return self.children[pos][ord]
        while True:
            nexts = self.children[pos]
            if nexts[ord] != -1:
                return nexts[ord]
            if pos == 0:
                return 0
            pos = self.link[pos]

    def clear(self) -> None:
        self.wordPos.clear()
        self.parent.clear()
        self.children.clear()
        self.link.clear()
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
        link = self.link
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
                from_, to = self.link[v], v
                arr1, arr2 = res[from_], res[to]
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
                arr3 += arr1[i:]
                arr3 += arr2[j:]
                res[to] = arr3
        return res

    def dp(self) -> Generator[Tuple[int, int], None, None]:
        """按照拓扑序进行转移, 每次返回(link, cur)."""
        link = self.link
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
        link = self.link
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


def min2(a: int, b: int) -> int:
    return a if a < b else b


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
                pos = acm.move(pos, ord(char))
                cur = pos
                while cur != 0:
                    dp[i + 1] = min2(dp[i + 1], dp[i + 1 - nodeDepth[cur]] + nodeCosts[cur])
                    cur = acm.link[cur]
            return dp[n] if dp[n] != INF else -1
