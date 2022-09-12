# https://www.acwing.com/problem/content/description/1287/

# 某人读论文，一篇论文是由许多单词组成的。
# 但他发现一个单词会在论文中出现很多次，现在他想知道每个单词分别在论文中出现多少次。
import string
from collections import defaultdict, deque, namedtuple, Counter
from typing import Iterable


# !给定k个单词和一段包含n个字符的文章,求有多少个单词在文章里`出现过`。
# 若使用KMP算法,则每个模式串T,都要与主串S进行一次匹配,
# 总时间复杂度为O(n×k+m),其中n为主串S的长度,m为各个模式串的长度之和,k为模式串的个数。
# !而采用AC自动机,时间复杂度只需O(n+m)。


def useAutoMaton(charset: Iterable[str] = string.ascii_lowercase):
    nodeId = 0
    trie = [defaultdict(int)]
    nexts = [0]  # kmp算法的nexts数组,失配指针
    count = [0]
    exists = [[]]

    def insert(pattern: str) -> None:
        nonlocal nodeId
        root = 0
        for char in pattern:
            if not trie[root][char]:
                nodeId += 1
                trie[root][char] = nodeId
                trie.append(defaultdict(int))
                nexts.append(0)
                count.append(0)
                exists.append([])
            root = trie[root][char]
        count[root] += 1
        exists[root].append(len(pattern))

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

    def search(word: str) -> Counter:
        """读入文章开始查询 `word` 中包含了AC自动机里的几个pattern"""
        res = Counter()
        root = 0
        for i, char in enumerate(word):
            cur = root = trie[root][char]
            while cur and exists[cur]:
                for len_ in exists[cur]:
                    res[word[i - len_ + 1 : i + 1]] += 1
                cur = nexts[cur]
        return res

    return namedtuple("AutoMaton", ["insert", "build", "search"])(insert, build, search)


n = int(input())
ac = useAutoMaton()
words = [""] * n

for i in range(n):
    words[i] = input()
    ac.insert(words[i])

ac.build()

res = ac.search(" ".join(words))
for word in words:
    print(res[word])
