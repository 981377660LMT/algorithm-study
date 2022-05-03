# https://www.acwing.com/problem/content/1284/

# 给定 n 个长度不超过 50 的由小写英文字母组成的单词，以及一篇长为 m 的文章。
# 请问，其中有多少个单词在文章中出现了。
# 注意：每个单词不论在文章中出现多少次，仅累计 1 次。
import sys
import string
from collections import defaultdict, deque, namedtuple
from typing import Iterable

input = lambda: sys.stdin.readline().strip()


def useAutoMaton(charset: Iterable[str]):
    nodeId = 0
    trie = [defaultdict(int)]
    next = [0]  # kmp算法的next数组,失配指针
    count = [0]

    def insert(word: str) -> None:
        nonlocal nodeId

        root = 0
        for char in word:
            if trie[root][char] == 0:
                trie[root][char] = nodeId
                nodeId += 1
                trie.append(defaultdict(int))
                next.append(0)
                count.append(0)
            root = trie[root][char]
        count[root] += 1

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
                if child == 0:
                    trie[cur][char] = trie[next[cur]][char]
                else:
                    # 孩子的失配指针指向父亲的失配指针的孩子 (四边形)
                    next[child] = trie[next[cur]][char]
                    queue.append(child)

    def search(pattern: str) -> int:
        """读入文章开始查询 `pattern` 中包含了AC自动机里的几个单词"""
        res = 0
        root = 0
        for char in pattern:
            root = trie[root][char]
            cur = root
            while cur and count[cur] != -1:
                res += count[cur]
                # 出现过就不用再往回跳了
                count[cur] = -1
                cur = next[cur]
        return res

    return namedtuple('AutoMaton', ['insert', 'build', 'search'])(insert, build, search)


if __name__ == '__main__':
    a = useAutoMaton(string.ascii_lowercase)
    a.insert('abc')
    a.insert('abc')
    a.insert('abc')
    a.insert('abcd')
    a.insert('abcde')
    a.insert('abcdef')
    a.insert('abcdefg')
    a.build()
    print(a.search('abcdefg'))
    # T = int(input())

    # for _ in range(T):
    #     autoMaton = useAutoMaton(string.ascii_lowercase)
    #     n = int(input())
    #     for _ in range(n):
    #         autoMaton.insert(input())

    #     autoMaton.build()

    #     pattern = input()
    #     print(autoMaton.search(pattern))

