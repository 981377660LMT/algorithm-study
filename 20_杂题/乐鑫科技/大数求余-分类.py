# 求一个 n (1 < n <= 36)进制正整数num除以 m (1 < m <= 36) 的余数。

radix, m = map(int, input().split())
chars = input()

i = 0
while i < len(chars):
    char = chars[i]
    if char.isdigit():
        if int(char) >= radix:
            break
    elif char.isalpha():
        if ord(char) - ord('a') + 10 >= radix:
            break
    i += 1

num = chars[:i]
res = 0
for j in range(i):
    char = num[j]
    cur = 0
    if char.isdigit():
        cur = int(char)
    elif char.isalpha():
        cur = ord(char) - ord('a') + 10
    # 进制转换
    res = res * radix + cur
    res %= m
print(res)

