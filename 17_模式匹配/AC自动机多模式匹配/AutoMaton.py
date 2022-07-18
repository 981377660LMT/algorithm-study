import string
from collections import defaultdict, deque, namedtuple
from typing import Iterable, List, Tuple


# 给定k个单词和一段包含n个字符的文章,求有多少个单词在文章里`出现过`。
# 若使用KMP算法,则每个模式串T,都要与主串S进行一次匹配,
# 总时间复杂度为O(n×k+m),其中n为主串S的长度,m为各个模式串的长度之和,k为模式串的个数。
# 而采用AC自动机,时间复杂度只需O(n+m)。


def useAutoMaton(charset: Iterable[str], maxLen: int):
    nodeId = 0
    trie = [defaultdict(int) for _ in range(maxLen)]
    nexts = [0] * maxLen  # kmp算法的nexts数组,失配指针
    count = [0] * maxLen
    exists = [[] for _ in range(maxLen)]

    def insert(word: str) -> None:
        nonlocal nodeId
        root = 0
        for char in word:
            if not trie[root][char]:
                nodeId += 1
                trie[root][char] = nodeId
            root = trie[root][char]
        count[root] += 1
        exists[root].append(len(word))

    def build() -> None:
        """bfs,字典树的每个结点添加失配指针,结点要跳转到哪里

        AC自动机的失配指针指向的节点代表的字符串是当前节点代表的字符串的最长后缀。
        不空,失配指针指；空,自己去指
        """
        queue = deque(trie[0].values())
        while queue:
            cur = queue.popleft()
            for char in charset:
                child = trie[cur][char]
                # 孩子指向失配指针的孩子 (三角形)
                if not child:
                    trie[cur][char] = trie[nexts[cur]][char]
                else:
                    # 孩子的失配指针指向父亲的失配指针的孩子 (四边形)
                    nexts[child] = trie[nexts[cur]][char]
                    queue.append(child)

    def search(pattern: str) -> List[Tuple[int, str]]:
        """读入文章开始查询 `pattern` 中包含了AC自动机里的几个单词"""
        res = []
        root = 0
        for i, char in enumerate(pattern):
            cur = root = trie[root][char]
            while cur and exists[cur]:
                for len_ in exists[cur]:
                    res.append((i - len_ + 1, pattern[i - len_ + 1 : i + 1]))
                cur = nexts[cur]
        return res

    return namedtuple("AutoMaton", ["insert", "build", "search"])(insert, build, search)


if __name__ == "__main__":
    autoMaton = useAutoMaton(string.ascii_lowercase, 10000)
    autoMaton.insert("he")
    autoMaton.insert("she")
    autoMaton.insert("his")
    autoMaton.insert("hers")
    autoMaton.build()
    assert autoMaton.search("ahishershe") == [
        (1, "his"),
        (3, "she"),
        (4, "he"),
        (4, "hers"),
        (7, "she"),
        (8, "he"),
    ]
