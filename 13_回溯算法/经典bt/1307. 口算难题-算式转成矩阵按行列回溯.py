from typing import List

# 2 <= words.length <= 5
# 1 <= words[i].length, results.length <= 7
# words[i], result 只含有大写英文字母
# 表达式中使用的不同字符数最大为 10


# 转成矩阵 dfs(row,col) 的形式回溯
# 1307. 口算难题-算式转成矩阵按行列回溯
class Solution:
    def isSolvable(self, words: List[str], result: str) -> bool:
        def bt(col: int, row: int, curSum: int) -> bool:
            if col == COL:
                return curSum == 0
            if row == ROW:
                return curSum % 10 == 0 and bt(col + 1, 0, curSum // 10)

            word = words[row]
            if col >= len(word):
                return bt(col, row + 1, curSum)

            char = word[~col]
            sign = 1 if row < ROW - 1 else -1
            if char in numByChar:
                # 注意此时特判前导0
                if numByChar[char] != 0 or (len(word) == 1 or col != len(word) - 1):
                    return bt(col, row + 1, curSum + sign * numByChar[char])
                return False
            else:
                cands = []
                for i in range(10):
                    if i in charByNum:
                        continue
                    # 注意此时特判前导0
                    if i == 0 and (len(word) > 1 and col == len(word) - 1):
                        continue
                    cands.append(i)

                for select in cands:
                    numByChar[char] = select
                    charByNum[select] = char
                    if bt(col, row + 1, curSum + sign * select):
                        return True
                    del charByNum[select]
                    del numByChar[char]

            return False

        words.append(result)  # 最后一行当成负号
        ROW, COL = len(words), max(map(len, words))
        numByChar = {}
        charByNum = {}
        return bt(0, 0, 0)


print(Solution().isSolvable(words=["SEND", "MORE"], result="MONEY"))
# 输入：words = ["SEND","MORE"], result = "MONEY"
# 输出：true
# 解释：映射 'S'-> 9, 'E'->5, 'N'->6, 'D'->7, 'M'->1, 'O'->0, 'R'->8, 'Y'->'2'
# 所以 "SEND" + "MORE" = "MONEY" ,  9567 + 1085 = 10652

