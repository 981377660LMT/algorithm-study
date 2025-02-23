# Binomial Coefficient

二项式系数/杨辉三角/组合数

https://oi-wiki.org/math/number-theory/lucas/

1. https://judge.yosupo.jp/problem/binomial_coefficient
   模不为素数的大数二项式系数 => 扩展卢卡斯定理
   数据范围：`1≤k≤n≤1e18 ，2≤mod≤1e6 ，不保证 mod 是质数。`
   预处理时间复杂度：`O(modlogmod)`，单次查询时间复杂度：`O(lognlogmod)`
2. https://judge.yosupo.jp/problem/binomial_coefficient_prime_mod
   模素数的二项式系数 => 求逆元
   数据范围：`1≤k≤n≤1e7 ，1≤mod≤2**30 ，mod 是质数。`
   预处理时间复杂度：`O(n)`，单次查询时间复杂度：`O(1)`
