# 给定k个单词和一段包含n个字符的文章,求有多少个单词在文章里`出现过`。
# 若使用KMP算法,则每个模式串T,都要与主串S进行一次匹配,
# !总时间复杂度为O(S1*k+S2),其中S1为主串S的长度,S2为`各个模式串的长度之和`,k为模式串的个数。
# !而采用AC自动机,时间复杂度只需O(S1+S2)。


from collections import deque
from typing import List, Tuple


class AhoCorasick:
    """https://ikatakos.com/pot/programming_algorithm/string_search"""

    __slots__ = ("_patterns", "_children", "_match", "_fail")

    def __init__(self, patterns: List[str]):
        self._patterns = patterns
        self._children = [{}]
        self._match = [[]]  # match[i] 表示节点i对应的字符串在patterns中的下标

        for pi, pattern in enumerate(patterns):
            if not pattern:
                continue
            self._insert(pi, pattern)

        self._fail = [0] * len(self._children)
        self._buildFail()

    def search(self, target: str) -> List[Tuple[int, int, int]]:
        """查询各个模式串在主串`target`中出现的`[起始索引,结束索引,模式串的索引]`"""
        match, patterns = self._match, self._patterns

        root = 0
        res = []
        for i, char in enumerate(target):
            root = self._next(root, char)
            res.extend((i - len(patterns[m]) + 1, i, m) for m in match[root])
        return res

    def _insert(self, pi: int, pattern: str) -> None:
        root = 0
        for char in pattern:
            if char in self._children[root]:
                root = self._children[root][char]
            else:
                len_ = len(self._children)
                self._children[root][char] = len_
                self._children.append({})
                self._match.append([])
                root = len_
        self._match[root].append(pi)

    def _buildFail(self) -> None:
        """bfs,字典树的每个结点添加失配指针,结点要跳转到哪里

        AC自动机的失配指针指向的节点所代表的字符串 是 当前节点所代表的字符串的 最长后缀。
        """
        children, match, fail = self._children, self._match, self._fail

        queue = deque(children[0].values())
        while queue:
            cur = queue.popleft()
            fafail = fail[cur]
            for char, child in children[cur].items():
                fail[child] = self._next(fafail, char)
                match[child].extend(match[fail[child]])
                queue.append(child)

    def _next(self, fafil: int, char: str) -> int:
        """沿着失配链,找到一个节点fafail,具有char的子节点"""
        while True:
            if char in self._children[fafil]:
                return self._children[fafil][char]
            if fafil == 0:
                return 0
            fafil = self._fail[fafil]


if __name__ == "__main__":

    class Solution:
        def multiSearch(self, big: str, smalls: List[str]) -> List[List[int]]:
            """多模式匹配indexOfAll"""
            ac = AhoCorasick(smalls)
            match = ac.search(big)
            res = [[] for _ in range(len(smalls))]
            for start, _, wordId in match:
                res[wordId].append(start)
            return res

    print(Solution().multiSearch("mississippi", ["is", "ppi", "hi", "sis", "i", "ssippi"]))
