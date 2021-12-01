# 如果 A = "abc" 且 B = "cde"，那么就有 'a' == 'c', 'b' == 'd', 'c' == 'e'。
# 利用 A 和 B 的等价信息，找出并返回 S 的按字典序排列最小的等价字符串。
class Solution:
    def smallestEquivalentString(self, s1: str, s2: str, baseStr: str) -> str:
        # 并查集。
        # 每次合并时保留字典序小的根即可。
        def union(x: str, y: str) -> None:
            px, py = sorted([find(x), find(y)])
            parent[py] = px

        def find(x: str) -> str:
            if parent.setdefault(x, x) != x:
                parent[x] = find(parent[x])
            return parent[x]

        parent = dict()
        for a, b in zip(s1, s2):
            union(a, b)
        return ''.join(find(char) for char in baseStr)


print(Solution().smallestEquivalentString(s1="parker", s2="morris", baseStr="parser"))
