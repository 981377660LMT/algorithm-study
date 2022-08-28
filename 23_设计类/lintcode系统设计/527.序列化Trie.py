# 序列化Trie、反序列化Trie


from collections import defaultdict


class TrieNode:
    def __init__(self):
        self.children = defaultdict(TrieNode)


class Solution:
    def serialize(self, root: "TrieNode") -> str:
        """序列化Trie"""

        def dfs(cur: "TrieNode") -> str:
            res = []
            for key, next in cur.children.items():
                res.append(key)
                res.append(dfs(next))
            return f'<{"".join(res)}>'  # 注意返回时才带上<>

        return dfs(root)

    def deserialize(self, data: str) -> "TrieNode":
        """反序列化Trie"""

        def dfs(cur: str) -> "TrieNode":
            res = TrieNode()
            depth = 0
            key, child = "", []

            for char in cur:
                if char == "<":
                    depth += 1
                    if depth >= 2:
                        child.append("<")
                elif char == ">":
                    depth -= 1
                    if depth >= 1:
                        child.append(">")
                    if depth == 1:
                        res.children[key] = dfs("".join(child))  # 子节点
                        key, child = "", []
                else:
                    if depth == 1:
                        key = char
                    elif depth >= 2:
                        child.append(char)

            return res

        return dfs(data)


if __name__ == "__main__":
    res = Solution().deserialize("<a<b<e<>>c<>d<f<>>>>")
    print(Solution().serialize(res))  # <a<b<>>>
