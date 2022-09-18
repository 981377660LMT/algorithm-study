1. 注意 splice 时是 mutable 返回值是删除的那些数
   如果要`切片删除元素拼接数组的话`
   可以先考虑 filter

```JS
nums.filter((_, i) => i !== removeIndex)
```

其次再考虑 slice

2. 注意不要边遍历数组边改变数组长度/删除这个数组的元素
