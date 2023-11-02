# 撤销背包

https://www.cnblogs.com/Schucking-Sattin/p/17726836.html

1. 撤销方式 1: 线段树分治(分治删点)，参考模版 `mutateWithoutOne`
   好处是不需要删除接口(undo)
2. 撤销方式 2: 直接撤销(naive)
   必须要有删除接口(undo)

---

可行性背包转方案数背包，以方便进行回退。

- 01 背包撤销操作：

```cpp
void insert(int x)
{
    for(int i = n; i >= x; --i)
        Inc(f[i], f[i - x]);
}
void remove(int x)
{
    for(int i = x; i <= n; ++i)
        Dec(f[i], f[i - x]);
}
```

- 完全背包撤销操作：

```cpp
void insert(int x)
{
    for(int i = x; i <= n; ++i)
        Inc(f[i], f[i - x]);
}
void remove(int x)
{
    for(int i = n; i >= x; --i)
        Dec(f[i], f[i - x]);
}
```

- 多重背包撤销操作
  https://leetcode.cn/circle/discuss/YnZBve/
