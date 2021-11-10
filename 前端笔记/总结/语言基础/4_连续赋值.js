var a = { n: 1 }
var b = a
a.x = a = { n: 2 } // 点的优先级大于等号的优先级

console.log(a.x)
console.log(b.x)
// 结果:
// undefined
// {n:2}

// 对象指针存在栈上 对象存在堆上
// 把 a.x = a = {n: 2}， 换成 b.x = a = {n: 2} 的时候，是不是会好理解了，虽然确实是这样。

// a 赋值，a 指向堆内存 {n:1}
// a = { n : 1 }
// b 赋值，b 也指向对内存 {n:1}
// b = a
// .的优先级大于=，所以优先赋值。ps：此时a.x已经绑定到了{n: 1 , x: undefined}被等待赋值
// a.x = undefined

// a // {n: 1 , x: undefined}
// b // 也指向上行的堆内存
// 同等级赋值运算从右到左，a改变堆内存指向地址，所以a = {n: 2},
// a.x = a = {n: 2};
// 因为a.x已经绑定到了{n: 1 , x: undefined}这个内存地址，所以相当于
// {n: 1 , x: undefined}.x = {n: 2}
