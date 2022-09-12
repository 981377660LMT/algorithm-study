# n<=100 字符串长度<=200
# !找到在密文各个数字串中都没有出现过的最短子串，
# !若有多个，选字典序最小的
# https://leetcode.cn/problems/lJDUTU/comments/


from itertools import product
from typing import Generator, List


n = int(input())
words = [input() for _ in range(n)]


def genString() -> Generator[str, None, None]:
    len_ = 1
    while True:
        for cur in product(range(10), repeat=len_):
            yield "".join(map(str, cur))
        len_ += 1


def solve(secrets: List[str]) -> str:
    """找到在密文words各个数字串中都没有出现过的最短子串,若有多个,选字典序最小的"""
    gen_ = genString()
    while True:
        cand = next(gen_)
        if all(cand not in word for word in secrets):
            print(cand)
            exit(0)


print(solve(words))

##############################################################################
# # 字典树+迭代加深
# import string
# from collections import defaultdict
# from typing import DefaultDict, List


# def dfs(root: DefaultDict[str, DefaultDict], depthLimit: int, path: List[str]) -> None:
#     global res
#     for char in string.digits:
#         path.append(char)
#         if char not in root:
#             res = "".join(path)
#             return
#         if len(path) < depthLimit:
#             dfs(root[char], depthLimit, path)
#         path.pop()


# n = int(input())
# trie = lambda: defaultdict(trie)
# trieRoot = trie()
# for i in range(n):
#     word = input()
#     for i in range(len(word)):
#         root = trieRoot
#         cur = word[i:]
#         for char in cur:
#             root = root[char]

# res, depth = "", 1
# while not res:
#     dfs(trieRoot, depth, [])
#     depth += 1
# print(res)
