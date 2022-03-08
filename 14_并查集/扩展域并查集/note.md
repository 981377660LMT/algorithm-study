种类并查集也叫扩展域并查集，用于检测`矛盾冲突`关系
不在一组的关系可以使用`种类并查集`来做
每个 i 预留 2 个编号，0-n-1 是自己人编号，n-2n-1 是不喜欢的人编号，然后就按并查集来做

`敌人的敌人是朋友`

```Python
uf = UnionFindArray(n * 2 + 2)
for cur, next in dislikes:
    if uf.isConnected(cur, next):
        return False
    uf.union(cur, next + n)
    uf.union(cur + n, next)
return True
```

<!-- 例题：
将罪犯们在`两座`监狱内重新分配，以求产生的冲突事件影响力都较小，从而保住自己的乌纱帽。假设只要处于同一监狱内的某两个罪犯间有仇恨，那么他们一定会在每年的某个时候发生摩擦。 -->

我们得到的信息，不是哪些人应该在相同的监狱，而是哪些人应该在不同的监狱
敌人的敌人就是朋友，2 和 4 是敌人，2 和 1 也是敌人，所以这里，`1 和 4 通过 2+n 这个元素间接地连接起来了`。这就是种类并查集工作的原理。

总结：
`如果 x,y 同类，则 uf.union(x,y),uf.union(x+n,y+n)`
`如果 x,y 异类，则 uf.union(x,y+n),uf.union(x+n,y)`
