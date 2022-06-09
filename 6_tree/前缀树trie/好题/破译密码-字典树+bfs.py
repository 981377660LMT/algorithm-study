# 有一种数字串的加密方法，原文是一个数字串，密文是若干个数字串。
# 已知解密的方法是:找到在密文各个数字串中都没有出现过的最短子串，
# 若有多个，选字典序最小的。现在给定密文，请你帮忙破译。
# 1 ≤ n ≤ 100
# 字符串长度不大于 200
# 注意到所有字串的前缀至多2e4个
# 因此dfs最多搜索5层
# !因此可以用迭代加深 指定长度 按字典序搜搜到给定长度就返回
# !也可以bfs(只要不爆空间的话)

from collections import defaultdict, deque
import string
from typing import DefaultDict


def bfs(root: DefaultDict[str, DefaultDict]) -> str:
    """在字典树中找到最短的没有出现的前缀 如果长度一样 字典序要最小"""
    queue = deque([(root, '')])
    while queue:
        node, path = queue.popleft()
        for char in string.digits:
            if char not in node:
                return path + char
            queue.append((node[char], path + char))


n = int(input())
trie = lambda: defaultdict(trie)
trieRoot = trie()
for i in range(n):
    word = input()
    # 加入所有子串的前缀
    for i in range(len(word)):
        root = trieRoot
        cur = word[i:]
        for char in cur:
            root = root[char]

print(bfs(trieRoot))

# 输入：
#      2
#      9527
#      0012345678
# 输出：02
