树上二分(Search(k)):在每个结点处,是向右走还是向左走

- **递归左子树->等待答案->递归右子树?** 的方法来查询 [l,r] 中最**右**的满足条件的点
- **递归右子树->等待答案->递归左子树?** 的方法来查询 [l,r] 中最**左**的满足条件的点

1. 值域线段树树上二分查询 remain 张卡片凑出来的`最大分数和`
   递归右子树->等待答案->递归左子树

```go
// 动态开点值域线段树写法
func (o *Node) Search(remain int) int {
	if o.left == o.right {
		return o.left * remain
	}
	o.pushDown()
	if o.rightChild.data.count >= remain {
		return o.rightChild.Search(remain)
	}
	return o.leftChild.Search(remain-o.rightChild.data.count) + o.rightChild.data.sum
}
```

```ts
// 普通值域线段树写法
function _search(root: number, left: number, right: number, remain: number): number {
  if (left === right) {
    return left * remain
  }

  this._pushDown(root)
  const mid = (left + right) >> 1
  if (this._count[(root << 1) | 1] >= remain) {
    return this._search((root << 1) | 1, mid + 1, right, remain)
  }
  return (
    this._search(root << 1, left, mid, remain - this._count[(root << 1) | 1]) +
    this._sum[(root << 1) | 1]
  )
}
```

2. 查找[0,n]内**第一个值大于 k**的数
   灵活地使用 atc 模板库的 minLeft/maxRight，有时候需要**反过来写**

```go
first := this.tree.MaxRight(0, func(e E) bool { return e.max < k }) // !找到第一个空座位>=k的行
first := this.tree.MaxRight(0, func(e E) bool { return e.sum == 0 }) // !找到第一个未坐满的行
```
