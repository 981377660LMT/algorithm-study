from typing import Generic, Iterable, List, TypeVar
from collections import defaultdict, deque


T = TypeVar("T", str, int)


class ACAutoMatonMapForMatchingLegacy(Generic[T]):
    """AC自动机,多模式串匹配."""

    __slots__ = ("matching", "wordCount", "children", "_fail", "_heavy")

    def __init__(self):
        self.matching = [[]]
        """trie树结点附带的信息.matching[i]表示节点(状态)i对应的字符串在patterns中的下标."""

        self.wordCount = []
        """trie树结点附带的信息.count[i]表示节点(状态)i的匹配个数."""

        self.children = [{}]
        """trie树.children[i]表示节点(状态)i的所有子节点,0 表示虚拟根节点."""

        self._fail = []
        """fail[i]表示节点(状态)i的失配指针."""

        self._heavy = False
        """构建时是否在matching中处理出每个结点匹配到的模式串id."""

    def insert(self, pid: int, pattern: Iterable[T]) -> None:
        """将模式串`pattern`插入到Trie树中.模式串一般是`被禁用的单词`.

        Args:
            pid (int): 模式串的唯一标识id.
            pattern (str): 模式串.
        """
        if not pattern:
            return
        children, matching = self.children, self.matching
        root = 0
        for char in pattern:
            nexts = children[root]
            if char in nexts:
                root = nexts[char]
            else:
                nextState = len(self.children)
                nexts[char] = nextState
                root = nextState
                children.append({})
                matching.append([])
        matching[root].append(pid)

    def build(self, heavy=False) -> None:
        """
        构建失配指针.
        bfs为字典树的每个结点添加失配指针,结点要跳转到哪里.
        AC自动机的失配指针指向的节点所代表的字符串 是 当前节点所代表的字符串的 最长后缀.

        Args:
            heavy (bool, optional): 是否处理出每个结点匹配到的模式串id. 默认为False.
        """
        self.wordCount = [len(m) for m in self.matching]
        self._fail = [0] * len(self.children)
        self._heavy = heavy
        children, matching, fail = self.children, self.matching, self._fail
        wordCount = self.wordCount
        queue = deque(children[0].values())
        while queue:
            cur = queue.popleft()
            curFail = fail[cur]
            for input_, next_ in children[cur].items():
                p = self.move(curFail, input_)
                fail[next_] = p  # !更新子节点的fail指针
                wordCount[next_] += wordCount[p]
                if heavy:
                    for m in matching[p]:  # !更新move状态匹配的模式串下标
                        matching[next_].append(m)
                queue.append(next_)

    def move(self, pos: int, input_: T) -> int:
        """
        从当前状态`pos`沿着字符`input_`转移到的下一个状态.
        沿着失配链上跳,找到第一个可由char转移的节点.
        """
        children, fail = self.children, self._fail
        while True:
            nexts = children[pos]
            if input_ in nexts:
                return nexts[input_]
            if pos == 0:
                return 0
            pos = fail[pos]

    def empty(self) -> bool:
        return len(self.children) == 1

    @property
    def size(self) -> int:
        return len(self.children)

    def __len__(self) -> int:
        return len(self.children)


if __name__ == "__main__":
    INF = int(1e18)

    def min2(a: int, b: int) -> int:
        return a if a < b else b

    class Solution:
        def minimumCost(self, target: str, words: List[str], costs: List[int]) -> int:
            book = defaultdict(lambda: INF)
            for w, c in zip(words, costs):
                book[w] = min2(book[w], c)
            words, costs = [], []
            for w, c in book.items():
                words.append(w)
                costs.append(c)

            ac = ACAutoMatonMapForMatchingLegacy()
            for i, w in enumerate(words):
                ac.insert(i, w)
            ac.build(heavy=True)

            matching = ac.matching
            n = len(target)
            dp = [0] + [INF] * n
            pos = 0
            for i, char in enumerate(target):
                pos = ac.move(pos, char)
                for wid in matching[pos]:
                    j = i - len(words[wid]) + 1
                    dp[i + 1] = min2(dp[i + 1], dp[j] + costs[wid])
            return dp[n] if dp[n] != INF else -1

    # 1032. 字符流
    # https://leetcode.cn/problems/stream-of-characters/
    class StreamChecker:
        def __init__(self, words: List[str]):
            self._ac = ACAutoMatonMapForMatchingLegacy()
            for i, word in enumerate(words):
                self._ac.insert(i, word)
            self._ac.build()
            self._state = 0

        def query(self, letter: str) -> bool:
            """从字符流中接收一个新字符，如果字符流中的任一非空后缀能匹配 words 中的某一字符串，返回 true"""
            self._state = self._ac.move(self._state, letter)
            return self._ac.wordCount[self._state] > 0
