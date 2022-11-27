# 高橋君は料理 1 から N の N 品の料理を作ろうとしています。
# 料理 i はオーブンを連続した Ti分間使うことで作れます。
# 2 つのオーブンを使えるとき、N 品の料理を全て作るまでに最短で何分かかりますか？
# n<=100
# Ti<=1000

# !注意有两个烤箱
# 因此最佳的策略是将料理分为相近的两半 取两者较大值
# 问题转换为01背包 dp[i][val] 表示前i个料理需要能否取到时间val


from typing import List


def cooking(times: List[int]) -> int:
    dp = set([0])
    for num in times:
        ndp = dp.copy()
        for pre in dp:
            ndp.add(pre + num)
        dp = ndp
    half = (sum(times) + 1) // 2
    while half not in dp:
        half += 1
    return half


if __name__ == "__main__":
    _ = int(input())
    costs = list(map(int, input().split()))
    print(cooking(costs))
