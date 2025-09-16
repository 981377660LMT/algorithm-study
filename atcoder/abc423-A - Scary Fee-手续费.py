# https://atcoder.jp/contests/abc423/tasks/abc423_a
# 手续费
# 你在银行账户有 X 日元。若一次提出 k×1000 日元（k 为非负整数），需要额外支付手续费 k×C 日元。总支出为 k×(1000+C) 日元，且不得使余额为负。
# 求可提出的最大金额（必须是 1000 的倍数）。


if __name__ == "__main__":
    X, C = map(int, input().split())
    k = X // (1000 + C)
    print(k * 1000)
