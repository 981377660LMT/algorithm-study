# C - Debug
# https://atcoder.jp/contests/abc394/tasks/abc394_c
#
# 只要字符串包含 WA 作为（连续的）子字符串，就重复以下操作:
# !将文本中出现的最前面的 WA 替换为 AC
# 注意，WWWWWWA 会变成 ACCCCCC
#
# !先替换最右边的也不会影响结果
# 只要字符串包含 WA 作为（连续的）子字符串，就重复执行以下操作。
# 将字符串中出现的 WA 中最后一个替换为 AC 。


s = input()
sb = list(s)
for i in range(len(sb) - 1, 0, -1):
    if sb[i - 1] == "W" and sb[i] == "A":
        sb[i - 1] = "A"
        sb[i] = "C"

print("".join(sb))
