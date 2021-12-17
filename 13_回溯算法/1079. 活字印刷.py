from itertools import permutations

# 1 <= tiles.length <= 7
# tiles 由大写英文字母组成
class Solution:
    def numTilePossibilities(self, tiles: str) -> int:
        res = set()
        for i in range(1, len(tiles) + 1):
            for j in permutations(tiles, i):
                res.add(j)
        return len(res)

    def numTilePossibilities2(self, tiles: str) -> int:
        def dfs(path, remain):
            if path != '':
                res.add(path)
            for i in range(len(remain)):
                dfs(path + remain[i], remain[:i] + remain[i + 1 :])

        res = set()
        dfs('', tiles)
        return len(res)


print(Solution().numTilePossibilities("AAB"))
# 输出：8
# 解释：可能的序列为 "A", "B", "AA", "AB", "BA", "AAB", "ABA", "BAA"。
