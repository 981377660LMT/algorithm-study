# E - Lucky 7 Battle(幸运7) 博弈dp


# 字符串 S 只包含数字，字符串 X 只包含 A 和 T，字符串 T 是空串。
# 当 Xi为 A 时，Aoki 操作，否则 Takahashi 操作。
# 每次操作可将 0 或 Si放在 T 末尾。
# !如果 T 是 7  的倍数则 Takahashi 获胜，否则 Aoki 获胜。


from functools import lru_cache
from typing import List


def lucky7Battle(digits: List[int], roles: str) -> bool:
    """
    - ゲームは N ラウンドからなり
    - Xi が A なら青木君が、T なら高橋君が以下の操作を行う
      操作:T の末尾に Si か 0 のどちらか一方を加える
    - #!T は 7 の倍数であれば高橋君の勝ちであり、そうでなければ青木君の勝ちです。

    Args:
        digits (str): 0,…,9 からなる長さ N の文字列
        roles (str): A,T からなる長さ N の文字列

    Returns:
        #!bool: 7 の倍数であるか
    """

    @lru_cache(None)
    def dfs(index: int, mod: int) -> bool:
        """当前在第index个位置,当前的mod值,takahashi是否能赢"""
        if index == n:
            return mod == 0

        role = roles[index]
        if role == "T":
            return dfs(index + 1, (mod * 10 + digits[index]) % 7) or dfs(index + 1, (mod * 10) % 7)
        return dfs(index + 1, (mod * 10 + digits[index]) % 7) and dfs(index + 1, (mod * 10) % 7)

    n = len(digits)
    res = dfs(0, 0)
    dfs.cache_clear()
    return res


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    digits = list(map(int, input()))
    roles = input()
    lucky7 = lucky7Battle(digits, roles)  # 7の倍数であるか
    print("Takahashi" if lucky7 else "Aoki")
