# 在依次输入单词 searchWord 的每一个字母后，
# 推荐 products 数组中前缀与 searchWord 相同的最多三个产品。
# 如果前缀相同的可推荐产品超过三个，请按字典序返回最小的三个。

from typing import List
from collections import defaultdict
from heapq import heappush, heappushpop


class WrappedStr(str):
    def __init__(self, string: str) -> None:
        super().__init__()
        self.string = string

    def __lt__(self, other: 'WrappedStr') -> bool:
        return self.string > other.string

    def __eq__(self, other: 'WrappedStr') -> bool:
        return self.string == other.string


# 如果是products不改变，则sorted(product)插入即可suggestion即可
# 如果是products可改变，则需要使用固定大小的堆来代替suggestion
class TrieNode:
    def __init__(self):
        self.children = defaultdict(TrieNode)
        self.suggestion = []

    def add_suggestion(self, product: str):
        if len(self.suggestion) < 3:
            heappush(self.suggestion, WrappedStr(product))
        else:
            heappushpop(self.suggestion, WrappedStr(product))

    def get_suggestion(self):
        return sorted(self.suggestion, reverse=True)


class Trie:
    def __init__(self) -> None:
        self.root = TrieNode()

    def insert(self, product: str):
        root = self.root
        for char in product:
            root = root.children[char]
            root.add_suggestion(product)

    def search(self, searchWord):
        res = []
        root = self.root
        for char in searchWord:
            root = root.children[char]
            res.append(root.get_suggestion())
        return res


class Solution:
    def suggestedProducts(self, products: List[str], searchWord: str) -> List[List[str]]:
        print(sorted(products))
        trie = Trie()
        for p in products:
            trie.insert(p)

        return trie.search(searchWord)


print(Solution().suggestedProducts(['mobile', 'mouse', 'moneypot', 'monitor', 'mousepad'], 'mouse'))

