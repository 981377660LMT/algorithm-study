# 给定k个单词和一段包含n个字符的文章,求有多少个单词在文章里`出现过`。
# 若使用KMP算法,则每个模式串T,都要与主串S进行一次匹配,
# !总时间复杂度为O(S1*k+S2),其中S1为主串S的长度,S2为`各个模式串的长度之和`,k为模式串的个数。
# !而采用AC自动机,时间复杂度只需O(S1+S2)。
# https://zhuanlan.zhihu.com/p/408665473
# https://ikatakos.com/pot/programming_algorithm/string_search
# AC自动机又叫AhoCorasick


from functools import lru_cache
from typing import Generator, Generic, Iterable, List, Tuple, TypeVar

INF = int(2e18)


T = TypeVar("T", str, int)


class ACAutoMatonMap(Generic[T]):
    """
    不调用 BuildSuffixLink 就是Trie, 调用 BuildSuffixLink 就是AC自动机.
    每个状态对应Trie中的一个结点, 也对应一个字符串.
    """

    __slots__ = ("wordPos", "_children", "_suffixLink", "_bfsOrder")

    def __init__(self):
        self.wordPos = []
        """wordPos[i] 表示加入的第i个模式串对应的节点编号."""
        self._children = [{}]
        """_children[v][c] 表示节点v通过字符c转移到的节点."""
        self._suffixLink = []
        """又叫fail.指向当前节点最长真后缀对应结点,例如"bc"是"abc"的最长真后缀."""
        self._bfsOrder = []
        """结点的拓扑序,0表示虚拟节点."""

    def addString(self, string: Iterable[T]) -> int:
        if not string:
            return 0
        pos = 0
        for char in string:
            nexts = self._children[pos]
            if char in nexts:
                pos = nexts[char]
            else:
                nextState = len(self._children)
                nexts[char] = nextState
                pos = nextState
                self._children.append({})
        self.wordPos.append(pos)
        return pos

    def addChar(self, pos: int, char: T) -> int:
        nexts = self._children[pos]
        if char in nexts:
            return nexts[char]
        nextState = len(self._children)
        nexts[char] = nextState
        self._children.append({})
        return nextState

    def move(self, pos: int, char: T) -> int:
        children, link = self._children, self._suffixLink
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
        self._suffixLink = [-1] * len(self._children)
        self._bfsOrder = [0] * len(self._children)
        head, tail = 0, 1
        while head < tail:
            v = self._bfsOrder[head]
            head += 1
            for char, next_ in self._children[v].items():
                self._bfsOrder[tail] = next_
                tail += 1
                f = self._suffixLink[v]
                while f != -1 and char not in self._children[f]:
                    f = self._suffixLink[f]
                self._suffixLink[next_] = f
                if f == -1:
                    self._suffixLink[next_] = 0
                else:
                    self._suffixLink[next_] = self._children[f][char]

    def getCounter(self) -> List[int]:
        """获取每个状态匹配到的模式串的个数."""
        counter = [0] * len(self._children)
        for pos in self.wordPos:
            counter[pos] += 1
        for v in self._bfsOrder:
            if v != 0:
                counter[v] += counter[self._suffixLink[v]]
        return counter

    def getIndexes(self) -> List[List[int]]:
        """获取每个状态匹配到的模式串的索引."""
        res = [[] for _ in range(len(self._children))]
        for i, pos in enumerate(self.wordPos):
            res[pos].append(i)
        for v in self._bfsOrder:
            if v != 0:
                from_, _children = self._suffixLink[v], v
                arr1, arr2 = res[from_], res[_children]
                arr3 = [0] * (len(arr1) + len(arr2))
                i, j, k = 0, 0, 0
                while i < len(arr1) and j < len(arr2):
                    if arr1[i] < arr2[j]:
                        arr3[k] = arr1[i]
                        i += 1
                    elif arr1[i] > arr2[j]:
                        arr3[k] = arr2[j]
                        j += 1
                    else:
                        arr3[k] = arr1[i]
                        i += 1
                        j += 1
                    k += 1
                while i < len(arr1):
                    arr3[k] = arr1[i]
                    i += 1
                    k += 1
                while j < len(arr2):
                    arr3[k] = arr2[j]
                    j += 1
                    k += 1
                res[_children] = arr3
        return res

    def dp(self) -> Generator[Tuple[int, int], None, None]:
        for v in self._bfsOrder:
            if v != 0:
                yield self._suffixLink[v], v

    def buildFailTree(self) -> List[List[int]]:
        adjList = [[] for _ in range(len(self._children))]
        for v in self._bfsOrder:
            if v != 0:
                adjList[self._suffixLink[v]].append(v)
        return adjList

    @property
    def size(self) -> int:
        return len(self._children)

    def __len__(self) -> int:
        return len(self._children)


if __name__ == "__main__":
    # 1032. 字符流
    # https://leetcode.cn/problems/stream-of-characters/description/
    class StreamChecker:
        __slots__ = ("ac", "counter", "pos")

        def __init__(self, wordPos: List[str]):
            self.ac = ACAutoMatonMap()
            for word in wordPos:
                self.ac.addString(word)
            self.ac.buildSuffixLink()
            self.counter = self.ac.getCounter()
            self.pos = 0

        def query(self, letter: str) -> bool:
            self.pos = self.ac.move(self.pos, letter)
            return self.counter[self.pos] > 0

    # https://leetcode.cn/problems/multi-search-lcci/
    # 给定一个较长字符串big和一个包含较短字符串的数组smalls，
    # 设计一个方法，根据smalls中的每一个较短字符串，对big进行搜索。
    # !输出smalls中的字符串在big里出现的所有位置positions，
    # 其中positions[i]为smalls[i]出现的所有位置。
    class Solution1:
        def multiSearch(self, big: str, smalls: List[str]) -> List[List[int]]:
            acm = ACAutoMatonMap()
            for s in smalls:
                acm.addString(s)
            acm.buildSuffixLink()

            indexes = acm.getIndexes()
            res = [[] for _ in range(len(smalls))]
            pos = 0
            for i, char in enumerate(big):
                pos = acm.move(pos, char)
                for index in indexes[pos]:
                    res[index].append(i - len(smalls[index]) + 1)
            return res

    # 2781. 最长合法子字符串的长度
    # https://leetcode.cn/problems/length-of-the-longest-valid-substring/
    # 给你一个字符串 word 和一个字符串数组 forbidden 。
    # 如果一个字符串不包含 forbidden 中的任何字符串，我们称这个字符串是 合法 的。
    # 请你返回字符串 word 的一个 最长合法子字符串 的长度。
    # 子字符串 指的是一个字符串中一段连续的字符，它可以为空。
    #
    # 1 <= word.length <= 1e5
    # word 只包含小写英文字母。
    # 1 <= forbidden.length <= 1e5
    # !1 <= forbidden[i].length <= 1e5
    # !sum(len(forbidden)) <= 1e7
    # forbidden[i] 只包含小写英文字母。
    #
    # 思路:
    # 类似字符流, 需要处理出每个位置为结束字符的包含至少一个模式串的`最短后缀`.
    # !那么此时左端点就对应这个位置+1
    class Solution:
        def longestValidSubstring(self, word: str, forbidden: List[str]) -> int:
            def min(a: int, b: int) -> int:
                return a if a < b else b

            def max(a: int, b: int) -> int:
                return a if a > b else b

            acm = ACAutoMatonMap()
            for s in forbidden:
                acm.addString(s)
            acm.buildSuffixLink()

            minBorder = [INF] * len(acm)
            for i, pos in enumerate(acm.wordPos):
                minBorder[pos] = min(minBorder[pos], len(forbidden[i]))
            for pre, cur in acm.dp():
                minBorder[cur] = min(minBorder[cur], minBorder[pre])

            res, left, pos = 0, 0, 0
            for right, char in enumerate(word):
                pos = acm.move(pos, char)
                left = max(left, right - minBorder[pos] + 2)
                res = max(res, right - left + 1)

            return res

    # https://www.luogu.com.cn/problem/P3311
    # 我们称一个正整数 x 是幸运数，当且仅当它的十进制表示中不包含数字串集合 words 中任意一个元素作为其子串。
    # ac自动机 + 数位dp
    def p3311() -> None:
        import sys

        sys.setrecursionlimit(int(1e6))
        input = lambda: sys.stdin.readline().rstrip("\r\n")
        MOD = int(1e9 + 7)

        upper = input()
        wordCount = int(input())
        words = [input() for _ in range(wordCount)]
        acm = ACAutoMatonMap[int]()
        for v in words:
            acm.addString((int(c) for c in v))
        acm.buildSuffixLink()

        nums = list(map(int, str(upper)))
        counter = acm.getCounter()

        @lru_cache(None)
        def dfs(index: int, hasLeadingZero: int, isLimit: bool, acPos: int) -> int:
            """当前在第index位,hasLeadingZero表示有前导0,isLimit表示是否贴合上界"""
            if index == len(nums):
                return int(not hasLeadingZero)

            res = 0
            up = nums[index] if isLimit else 9
            for cur in range(up + 1):
                if hasLeadingZero and cur == 0:
                    res += dfs(index + 1, True, (isLimit and cur == up), acPos)
                else:
                    nextPos = acm.move(acPos, cur)
                    if counter[nextPos] == 0:
                        res += dfs(index + 1, False, (isLimit and cur == up), nextPos)
            return res % MOD

        res = dfs(0, True, True, 0)
        dfs.cache_clear()
        print(res)

    p3311()
