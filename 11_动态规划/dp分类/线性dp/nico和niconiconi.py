# 给定一个字符串 S ，定义三种有价值的字符串： A = "nico" ，B = "niconi" , C = "niconiconi"
# 其中，字符串 A 的价值为 a， 字符串 B 的价值为 b，字符串 C 的价值为 c

# 第一行四个正整数 n,a,b,c 。 (1≤n≤300000，1≤a,b,c≤10^9)
# 第二行是一个长度为 n  的字符串。
n, a, b, c = list(map(int, input().split()))
s = input()
str_a = 'nico'
str_b = 'niconi'
str_c = 'niconiconi'

dp = [0] * n


for i in range(1, n):
    dp[i] = dp[i - 1]
    if i >= 3 and s[i - 3 : i + 1] == str_a:
        dp[i] = max(dp[i], dp[i - 4] + a)
    if i >= 6 and s[i - 5 : i + 1] == str_b:
        dp[i] = max(dp[i], dp[i - 6] + b)
    if i >= 9 and s[i - 9 : i + 1] == str_c:
        dp[i] = max(dp[i], dp[i - 10] + c)

print(dp[n - 1])
