# Kahan 求和

https://oi-wiki.org/misc/kahan-summation/

Kahan 求和 算法，又名补偿求和或进位求和算法，
是一个用来 降低有限精度浮点数序列累加值误差 的算法。
它主要通过保持一个单独变量用来累积误差（常用变量名为 c）来完成的。

在浮点加法计算中，交换律（commutativity）成立，但结合律（associativity）不成立。
也就是说，$a+b = b+a$ 但 $(a+b)+c \neq a+(b+c)$。
