python 里的 product 函数用于求`各个 iterable 对象的笛卡尔积`，非常适合`穷举所有对象的状态的组合`
比如这道题，每个棋子的下一步都有很多个状态 states (这个 states 只要是一个 iterable 的对象就行)，那么 product(states1,states2,states3,...)可以穷举出所有棋子状态的所有组合

一般 product 会配合数组解构使用，例如` product(*(iterable for _ in range(n)))`

枚举的解法：

1. 描述状态：起点，方向，距离 `(x,y,dx,dy,dist)`
2. 枚举方向+枚举步数
