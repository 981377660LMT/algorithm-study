# 给定一些括号字符串
# !任意排序能否形成一个合法的括号序列

from typing import List


def bracketSequencing(words: List[str]) -> bool:
    bad, good = [], []
    for word in words:
        cost, cashback = 0, 0  # 启动每个项目的资金 完成项目的回报
        for char in word:
            if char == "(":
                cashback += 1
            elif cashback > 0:
                cashback -= 1
            else:
                cost += 1
        if cost > cashback:
            bad.append((cost, cashback))
        else:
            good.append((cost, cashback))

    bad.sort(key=lambda x: x[1], reverse=True)  # !前面cashback越大就越容易
    good.sort(key=lambda x: x[0])  # !前面cost越小就越容易
    nums = good + bad

    cur = 0
    for cost, cashback in nums:
        if cur < cost:
            return False
        cur -= cost
        cur += cashback
    return cur == 0


n = int(input())
words = [input() for _ in range(n)]
print("Yes" if bracketSequencing(words) else "No")
