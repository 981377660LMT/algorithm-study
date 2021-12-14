from string import ascii_lowercase

# 我们从一块字母板上的位置 (0, 0) 出发，该坐标对应的字符为 board[0][0]。
# 用最小的行动次数让答案和目标 target 相同。你可以返回任何达成目标的路径。


# 1.初始化位置
# 2. 注意要先向左/上移动，再向右/下移动 因为右下方缺了一块


class Solution:
    def alphabetBoardPath(self, target: str) -> str:
        pos = {char: (i // 5, i % 5) for i, char in enumerate(ascii_lowercase)}
        row, col = 0, 0
        res = []
        for char in target:
            x, y = pos[char]
            if y < col:
                res.append('L' * (col - y))
            if x < row:
                res.append('U' * (row - x))
            if x > row:
                res.append('D' * (x - row))
            if y > col:
                res.append('R' * (y - col))

            res.append('!')
            row, col = x, y

        return ''.join(res)


print(Solution().alphabetBoardPath(target="leet"))
# 输出："DDR!UURRR!!DDD!"
