# B - 银行利息1%
# 浮点数取整
# !浮点数乘法取整 => 先将浮点数变为整数(*101),再取整(//100)

# 高橋くんはAtCoder銀行に 100 円を預けています。
# AtCoder銀行では、毎年預金額の 1 % の利子がつきます(複利、小数点以下切り捨て)。
# 利子以外の要因で預金額が変化することはないと仮定したとき、高橋くんの預金額が初めて X 円以上になるのは何年後でしょうか。
if __name__ == "__main__":
    target = int(input())
    money = 100
    for i in range(int(1e9)):
        if money >= target:
            print(i)
            break
        money = money * 101 // 100  # 注意不能用int(cur * 1.01) 因为浮点数乘法精度问题
