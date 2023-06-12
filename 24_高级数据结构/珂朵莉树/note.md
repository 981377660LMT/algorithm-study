只要一个支持快速
`Insert/Erase/Prev/Next` 的数据结构就能实现珂朵莉树

- map?
- fastset(64-ary tree)?
- linked list?

珂朵莉树的特点是，每个区间`块的值是相等的`，`值相等的区间可以自动合并`
可以配合线段树+遍历块来做一些半群做不到的修改

---

珂朵莉树，又叫基于数据随机的颜色段均摊，用于管理区间。它可以高效地进行区间的赋值、查询、遍历和删除操作。珂朵莉树的特点是，`维护的每段区间的值是相等的，值相等的区间可以自动合并。`
这里用到的珂朵莉树是具有以下三个接口的数据结构。

```ts
set: set(start: number, end: number, value: S): void，用于设置指定范围内的值。它接受起始位置 start、结束位置 end 和值 value 作为参数，并将该范围内的值设置为指定的值。

get: get(x: number, erase = false): [start: number, end: number, value: S] | undefined，用于获取包含特定位置 x 的区间信息。它返回一个包含起始位置 start、结束位置 end 和对应的值 value 的元组。如果 erase 参数为 true，则在获取值的同时会将该区间删除。

enumerateRange: enumerateRange(start: number, end: number, f: (start: number, end: number, value: S) => void, erase = false): void，用于遍历指定范围内的所有区间，并对每个区间执行回调函数 f。回调函数 f 接受每个区间的起始位置 start、结束位置 end 和对应的值 value 作为参数。如果 erase 参数为 true，则在遍历区间的同时会将遍历到的区间删除。
```

只要一个数据结构可以支持**快速寻找前驱和后继**以及**快速删除和插入元素**，那么就可以用来实现珂朵莉树，例如 `SortedDict/64 叉树/Van Emde Boas Tree/链表`，等等...

- 珂朵莉树的 set 操作推平区间，有利于降低复杂度
- 珂朵莉树不擅长范围查询(因为复杂度与范围内的区间个数成正比)，一般需要配合其他数据结构辅助范围查询

---

优化:如果区间染色是某个值时,js 可以考虑使用 Int8Array 或者 Uint32Array 上的 `fill 方法`填充
(类型数组的 fill 方法使用了 `memset`)
