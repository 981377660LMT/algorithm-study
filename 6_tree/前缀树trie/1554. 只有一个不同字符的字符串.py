#  !给定一个字符串列表 dict ，其中所有字符串的长度都`相同`。
#  当存在两个字符串在相同索引处只有一个字符不同时，返回 True ，否则返回 False 。
#  dict 中的字符数小于或等于 10^5 。
#  你可以以 O(n*m) 的复杂度解决问题吗？其中 n 是列表 dict 的长度，m 是字符串的长度。


# 将单词插入前缀树之前判断是否可以满足只有一个字符不同的条件
from typing import List
from Trie import Trie, TrieNode


class Solution:
    def differByOne(self, dict: List[str]) -> bool:
        def search(word: str) -> bool:
            def dfs(cur: "TrieNode", index: int, changed: bool) -> bool:
                if index == len(word):
                    return changed
                res = False
                for child in cur.children:
                    if child == word[index]:
                        res |= dfs(cur.children[child], index + 1, changed)
                    elif not changed:
                        res |= dfs(cur.children[child], index + 1, True)
                return res

            return dfs(trie.root, 0, False)

        trie = Trie()
        for word in dict:
            if search(word):
                return True
            trie.insert(word)
        return False


assert Solution().differByOne(["abcd", "acbd", "aacd"])
