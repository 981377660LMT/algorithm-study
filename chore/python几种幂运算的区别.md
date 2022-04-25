# 6042. 统计圆内格点数目

总结:`幂运算用 a*a 比pow(a,2)/a**2快`

`a**2和pow(a,2)`：求准确值，因为要支持大数 `logn`
`乘法a*a`:O(1)

![python几种幂运算的区别](image/python几种幂运算的区别/1650781677843.png)
![pow(a,b)和**调用的是同一个函数](image/python几种幂运算的区别/1650781689227.png)
![**超时](image/python几种幂运算的区别/1650781736265.png)
![*不超时](image/python几种幂运算的区别/1650781744476.png)

https://leetcode-cn.com/problems/count-lattice-points-inside-a-circle/solution/python-bao-li-jie-fa-chao-shi-de-yuan-yi-jwf7/

[Speed of calculating powers (in python)](https://stackoverflow.com/questions/1019740/speed-of-calculating-powers-in-python)
