class Solution:
    def maxCollectedFruits(self, fruits: List[List[int]]) -> int:
        n = len(fruits)
        total_fruits = 0
        visited = [[False] * n for _ in range(n)]

        # Kid 1 path (Main Diagonal)
        for i in range(n):
            if not visited[i][i]:
                total_fruits += fruits[i][i]
                visited[i][i] = True

        # Kid 2 path (Zigzag Path)
        i, j = 0, n - 1
        while i < n - 1:
            if not visited[i][j]:
                total_fruits += fruits[i][j]
                visited[i][j] = True
            if j > 0:
                i += 1
                j -= 1
            else:
                i += 1
        if not visited[n - 1][n - 1]:
            total_fruits += fruits[n - 1][n - 1]
            visited[n - 1][n - 1] = True

        # Kid 3 path (Bottom Row)
        for j in range(n):
            i = n - 1
            if not visited[i][j]:
                total_fruits += fruits[i][j]
                visited[i][j] = True

        return total_fruits
