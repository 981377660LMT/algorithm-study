多路归并的特征:
**寻找第 k 大/小的某个数**
其中一般 k<=2000

一般的流程

1. pq push 起始值 (val,row,col) 或者 (val,i,j)
   `这里多路指的是以哪个维度区分不同的路`
   例如:所在数组行/以哪个元素结尾/最后一个元素是第几个元素
2. 获取 k 个结果(push/pop)
3. push 下一个

```JS
 pq.push(起始值)
 while (pq.length && res.length < k) {
    const [_, i, j] = pq.shift()!
    res.push([nums1[i], nums2[j]])
    // ...
    pushNext(val,i,j+1)
  }
 return res
```
