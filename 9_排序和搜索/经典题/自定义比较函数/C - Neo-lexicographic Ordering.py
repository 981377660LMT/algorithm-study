# 自定义火星词典排序


from typing import List


def solve(s: str, words: List[str]) -> List[str]:
    order = {char: i for i, char in enumerate(s, 1)}
    words.sort(key=lambda word: [order[char] for char in word])
    return words


if __name__ == "__main__":
    s = input()  # 火星词典 zyxwvutsrqponmlkjihgfedcba
    n = int(input())
    words = [input() for _ in range(n)]
    print(*solve(s, words), sep="\n")
