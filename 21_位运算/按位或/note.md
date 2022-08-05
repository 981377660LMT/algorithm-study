## 与、或、异或 vs 两数和 三个重要的公式

不用加减乘除实现加法

1.  `(a & b) << 1 + (a ^ b) = a + b`

2.  `(a & b) + (a | b) = a + b`

3.  `(a | b).bit_count() + (a & b).bit_count() = a.bit_count() + b.bit_count()`

[与运算与两数和](../%E6%8C%89%E4%BD%8D%E4%B8%8E/D%20-%20AND%20and%20SUM.py)
[与运算与或运算](6127.%20%E4%BC%98%E8%B4%A8%E6%95%B0%E5%AF%B9%E7%9A%84%E6%95%B0%E7%9B%AE-%E8%84%91%E7%AD%8B%E6%80%A5%E8%BD%AC%E5%BC%AF.py)
