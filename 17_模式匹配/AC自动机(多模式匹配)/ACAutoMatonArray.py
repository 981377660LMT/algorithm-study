# 给定k个单词和一段包含n个字符的文章,求有多少个单词在文章里`出现过`。
# 若使用KMP算法,则每个模式串T,都要与主串S进行一次匹配,
# !总时间复杂度为O(S1*k+S2),其中S1为主串S的长度,S2为`各个模式串的长度之和`,k为模式串的个数。
# !而采用AC自动机,时间复杂度只需O(S1+S2)。
# https://zhuanlan.zhihu.com/p/408665473
# https://ikatakos.com/pot/programming_algorithm/string_search
# AC自动机又叫AhoCorasick


from collections import defaultdict, deque
from typing import Callable, DefaultDict, List, Optional


class ACAutoMatonArray:
    """AC自动机,多模式串匹配."""

    __slots__ = ("_patterns", "_children", "_wordCount", "_matching", "_fail", "_heavy")

    def __init__(self):
        self._patterns = []
        """模式串列表."""

        self._children = []
        """trie树.children[i]表示节点(状态)i的所有子节点,0 表示虚拟根节点."""

        self._wordCount = []
        """trie树结点附带的信息.count[i]表示节点(状态)i的匹配个数."""

        self._matching = [[]]
        """trie树结点附带的信息.matching[i]表示节点(状态)i对应的字符串在patterns中的下标."""

        self._fail = []
        """fail[i]表示节点(状态)i的失配指针."""

        self._heavy = False
        """构建时是否在matching中处理出每个结点匹配到的模式串id."""

    def insert(
        self, pid: int, pattern: str, didInsert: Optional[Callable[[int], None]] = None
    ) -> "ACAutoMatonArray":
        """将模式串`pattern`插入到Trie树中.模式串一般是`被禁用的单词`.

        Args:
            pid (int): 模式串的唯一标识id.
            pattern (str): 模式串.
            didInsert (Optional[Callable[[int], None]]): 模式串插入后的回调函数,入参为结束字符所在的结点(状态).
        """
        if not pattern:
            return self
        children, matching, patterns = self._children, self._matching, self._patterns
        root = 0
        for char in pattern:
            nexts = children[root]
            if char in nexts:
                root = nexts[char]
            else:
                nextState = len(self._children)
                nexts[char] = nextState
                root = nextState
                children.append({})
                matching.append([])
        matching[root].append(pid)
        patterns.append(pattern)
        if didInsert is not None:
            didInsert(root)
        return self

    def build(self, heavy=False, dp: Optional[Callable[[int, int], None]] = None) -> None:
        """
        构建失配指针.
        bfs为字典树的每个结点添加失配指针,结点要跳转到哪里.
        AC自动机的失配指针指向的节点所代表的字符串 是 当前节点所代表的字符串的 最长后缀.

        Args:
            heavy (bool, optional): 是否处理出每个结点匹配到的模式串id. 默认为False.
            dp (Optional[Callable[[int, int], None]], optional):
            AC自动机构建过程中的回调函数,入参为`(next_结点的fail指针, next_结点)`.
        """
        self._wordCount = [len(m) for m in self._matching]
        self._fail = [0] * len(self._children)
        self._heavy = heavy
        children, worCount, matching, fail = (
            self._children,
            self._wordCount,
            self._matching,
            self._fail,
        )
        queue = deque(children[0].values())
        while queue:
            cur = queue.popleft()
            curFail = fail[cur]
            for input_, next_ in children[cur].items():
                move = self.move(curFail, input_)
                fail[next_] = move  # !更新子节点的fail指针
                worCount[next_] += worCount[move]
                if heavy:
                    for m in matching[move]:  # !更新move状态匹配的模式串下标
                        matching[next_].append(m)
                if dp is not None:
                    dp(move, next_)
                queue.append(next_)

    def move(self, pos: int, input_: str) -> int:
        """
        从当前状态`pos`沿着字符`input_`转移到的下一个状态.
        沿着失配链上跳,找到第一个可由char转移的节点.
        """
        children, fail = self._children, self._fail
        while True:
            nexts = children[pos]
            if input_ in nexts:
                return nexts[input_]
            if pos == 0:
                return 0
            pos = fail[pos]

    def match(self, pos: int, s: str) -> DefaultDict[int, List[int]]:
        """从状态`pos`开始匹配字符串`s`.

        Args:
            pos (int): ac自动机的状态.根节点状态为0.
            s (str): 待匹配的字符串.

        Returns:
            DefaultDict[int, List[int]]: 每个模式串在`s`中出现的下标.
        """
        assert self._heavy, "需要调用build(heavy=True)构建AC自动机"
        matching, patterns = self._matching, self._patterns
        res = defaultdict(list)
        root = pos
        for i, char in enumerate(s):
            root = self.move(root, char)
            for m in matching[root]:
                res[m].append(i - len(patterns[m]) + 1)
        return res

    def count(self, pos: int) -> int:
        """当前状态`pos`匹配到的模式串个数."""
        return self._wordCount[pos]

    def accept(self, pos: int) -> bool:
        """当前状态`pos`是否为匹配状态."""
        return self._wordCount[pos] > 0


if __name__ == "__main__":

    def demo() -> None:
        ac = ACAutoMatonArray()
        ac.insert(0, "he")
        ac.insert(1, "she")
        ac.insert(2, "his")
        ac.insert(3, "hers")
        ac.build(heavy=True)
        print(ac.count(0))
        s = ac.move(0, "h")
        print(s)
        s = ac.move(s, "e")
        print(s)
        s = ac.move(s, "r")
        print(s)
        s = ac.move(s, "s")
        print(ac.match(s, "ushers"))

    demo()

    # https://leetcode.cn/problems/multi-search-lcci/
    class Solution:
        def multiSearch(self, big: str, smalls: List[str]) -> List[List[int]]:
            """多模式匹配indexOfAll"""
            ac = ACAutoMatonArray()
            for i, small in enumerate(smalls):
                ac.insert(i, small)
            ac.build(heavy=True)

            matching = ac.match(0, big)
            res = [[] for _ in range(len(smalls))]
            for wordId, starts in matching.items():
                for start in starts:
                    res[wordId].append(start)
            return res

    print(Solution().multiSearch("mississippi", ["is", "ppi", "hi", "sis", "i", "ssippi"]))

    # 1032. 字符流
    # https://leetcode.cn/problems/stream-of-characters/
    class StreamChecker:
        def __init__(self, words: List[str]):
            self._ac = ACAutoMatonArray()
            for i, word in enumerate(words):
                self._ac.insert(i, word)
            self._ac.build()
            self._state = 0

        def query(self, letter: str) -> bool:
            """从字符流中接收一个新字符，如果字符流中的任一非空后缀能匹配 words 中的某一字符串，返回 true"""
            self._state = self._ac.move(self._state, letter)
            return self._ac.accept(self._state)

    # Your StreamChecker object will be instantiated and called as such:
    # obj = StreamChecker(words)
    # param_1 = obj.query(letter)
