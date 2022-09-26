# 定义LCP(X,Y)为字符串X,Y的公共前缀长度(LCP)。
# 给定长度为N的字符串S，设S表示从第i个字符开始的S的后缀(就是后缀数组里的那些后缀)。
# !计算出:对于k=1,2,...,N,LCP(Sk,S1)+LCP(Sk,S2)+ +...+LCP(Sk,SN)的值。
# n<=1e6
# https://atcoder.jp/contests/abc213/tasks/abc213_f

# !即求每个后缀与所有后缀的公共前缀长度和。
# 任意两后缀的 LCP
# !lcp(sa[i], sa[j])=min{height[i＋1..j}
# !单调栈 看每个元素的影响范围
