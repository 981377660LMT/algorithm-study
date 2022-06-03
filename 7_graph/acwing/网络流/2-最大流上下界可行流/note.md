https://www.acwing.com/solution/content/24047/
每个流量有个下界，这是不同于最基本的模型（基本模型下界为 0）的一点，我们的思路是进行转换，从而让下界都变为 0。
https://www.acwing.com/solution/content/68545/

**多源汇最大流**
新建超级源点 S,向所有源点连容量为+∞ 的边即可。汇点同理。问题即转化为普通最大流问题。
