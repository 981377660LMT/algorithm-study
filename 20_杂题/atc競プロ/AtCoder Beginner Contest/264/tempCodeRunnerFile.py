r _ in range(2)] for _ in range(2)]
dp[0][0][0] = 0
dp[1][0][0] = rowCost[0]
dp[0][1][0] = colCost[0]
dp[1][1][0] = rowCost[0] + colCost[0]

for r in range(ROW):
    for c in range(COL):
        pos = r * COL + c
        for rowFlip in range(2):
            for colFlip in range(2):
                if r + 1 < ROW:
                    if grid[r + 1][c] == (grid[r][c] ^ rowFlip):
                        dp[0][colFlip][pos + COL] = min(
                            dp[0][colFlip][pos + COL], dp[rowFlip][colFlip][pos]
                        )
                    else:
                        dp[1][colFlip][pos + COL] = min(
                            dp[1][colFlip][pos + COL], dp[rowFlip][colFlip][pos] + rowCost[r + 1]
                        )
                if c + 1 < COL:
                    if grid[r][c + 1] == (grid[r][c] ^ colFlip):
                        dp[rowFlip][0][pos + 1] = min(
                            dp[rowFlip][0][pos + 1], dp[rowFlip][colFlip][pos]
                        )
                    else:
                        dp[rowFlip][1][pos + 1] = min(
                            dp[rowFlip][1][pos + 1], dp[rowFlip][colFlip][pos] + colCost[c + 1]
                        )

print(min(dp[1][1][-1], dp[1][0][-