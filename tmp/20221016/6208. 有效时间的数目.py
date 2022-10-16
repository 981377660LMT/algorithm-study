# 6208. 有效时间的数目


MOD = int(1e9 + 7)
INF = int(1e20)

# 在字符串 time 中，被字符 ? 替换掉的数位是 未知的 ，被替换的数字可能是 0 到 9 中的任何一个。

# 请你返回一个整数 answer ，将每一个 ? 都用 0 到 9 中一个数字替换后，可以得到的有效时间的数目。
# "00" <= hh <= "23"
# "00" <= mm <= "59"


class Solution:
    def countTime(self, time: str) -> int:
        """dfs每位避免了复杂的处理"""

        def check(s: str) -> bool:
            hh, mm = s.split(":")
            return 0 <= int(hh) <= 23 and 0 <= int(mm) <= 59

        def dfs(index: int, path: str) -> None:
            if len(path) == 5:
                self.res += check(path)
                return
            if time[index] != "?":
                dfs(index + 1, path + time[index])
            else:
                for i in range(10):
                    dfs(index + 1, path + str(i))

        self.res = 0
        dfs(0, "")
        return self.res
