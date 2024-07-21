https://qiita.com/ktateish/items/ab2df3e0864d2e931bf2

- AtCoder のジャッジシステムは 64bit 環境であり int 型が 64bit サイズになっているため、int64 型をわざわざ使い分けなくて良いのは楽です。
- slice の先頭と末尾からであれば要素を O(1) で削除できます

## golang 切片

1. 切片的[:]操作取的是一段 subarray，会共享原切片底层数组(这一点与 python 不同)
2. append 进行扩容时，什么时候会共享原切片底层数组？什么时候会使用新数组来创建切片？
   append 之后元素个数**小于等于原切片剩余容量时，会共享原切片底层数组，否则会扩容使用新数组来创建切片**
3. 扩容系数
   长度小于 1024 时，扩容系数为 2，长度大于等于 1024 时，扩容系数为 1.25
4. 如何对切片进行浅拷贝?

   > 浅拷贝指的是，把变量里面存的内存地址拷贝了，所指向的真实值并没拷贝。

   第一种方法是 append+解构 (类似 js 的 [...nums] )

   ```go
   func main() {
   	slice1 := []int{1, 2, 3, 4, 5}
   	sliceCopy := append([]int{}, slice1...)

    // !不需要写类型拷贝
    // sliceCopy := append(slice1[:0:0], slice1...)
    // 注意容量要为0，否则会共享底层数组

   	sliceCopy[0] = 100
   	fmt.Println(slice1, sliceCopy) // [1 2 3 4 5] [100 2 3 4 5]
   }
   ```

   第二种方法是 copy

   ```go
   func main() {
   	// 二维数组拷贝
   	slice2 := [][]int{{1, 2, 3}, {4, 5, 6}}
   	sliceCopy := make([][]int, len(slice2))
   	for i := range slice2 {
   		sliceCopy[i] = make([]int, len(slice2[i]))
   		copy(sliceCopy[i], slice2[i])
   	}
   	sliceCopy[0][0] = 100
   	fmt.Println(slice2, sliceCopy) // [[1 2 3] [4 5 6]] [[100 2 3] [4 5 6]]
   }
   ```

5. 切片内存储类似 tuple 的对时，可以用结构体实现
6. 切片初始化
   利用切片声明一个`静态数组`

   ```go
   nums := make([]int, 5)
   类似js的 const nums = Array(5).fill(0)
   ```

   声明一个`动态数组`(知道容量的话最好写上容量)

   ```go
   nums := make([]int, 0, 5)
   类似js的 const nums = []
   ```

7. make 和 new 的区别
   make 用于 chan map slice 的初始化，new 用于类型的属性初始化
   他们都返回的是指针(`注意 chan/map/slice 本身就是指针`)
8. make 初始化 struct 是`深层次的初始化`

   ```go
   package main

   import "fmt"

   type Fish struct {
   	id   int
   	name string
   	Metadata
   }

   type Metadata struct {
   	age   int
   	speed int
   	pos   float64
   }

   func main() {
   	fish := make([]Fish, 5)
   	fmt.Println(fish) // [{0  {0 0 0}} {0  {0 0 0}} {0  {0 0 0}} {0  {0 0 0}} {0  {0 0 0}}]
   }
   ```

9. 何时需要传递 slice 的指针?

- 如果只需要修改某个元素的值，那么传递 slice 指针是不必要的，因为 slice 本身就是指针
- 但是如果需要修改 slice 的 begin/len/cap，那么就需要传递 slice 指针了，
  避免扩容生成新的 slice 而无法修改原 slice

```go
func main() {
 	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
 	mut(&nums)
 	fmt.Println(nums)
 }

 func mut(nums *[]int) {
 	(*nums)[3] = 100          // 这种情况没有必要用指针
 	*nums = append(*nums, 11) // 这种情况必须用指针
 	*nums = (*nums)[1:5]      // 这种情况必须用指针
 }
```

## golang map

1. Creating map with/without make
   https://stackoverflow.com/questions/16959992/creating-map-with-without-make
   https://stackoverflow.com/questions/31064688/which-is-the-nicer-way-to-initialize-a-map
   **One allows you to initialize capacity, one allows you to initialize values:**

   ```go
    // make
    mp1 := make(map[string]int, 1024)
    // literal
    mp2 := map[string]int{
        "foo": 1,
        "bar": 2,
    }
   ```

2. 记得初始化 map

## 封装模板的方法

1. struct + method

## golang 浅拷贝与深拷贝

https://blog.csdn.net/weixin_40165163/article/details/90680466
https://github.com/981377660LMT/ts/issues/136

## for 语句

https://blog.csdn.net/jfkidear/article/details/89813758

for 语句的功能用来指定重复执行的语句块，for 语句中的表达式有三种：
官方的规范： ForStmt = "for" [ Condition | ForClause | RangeClause ] Block .

1. Condition = Expression .
2. ForClause = [ InitStmt ] “;” [ Condition ] “;” [ PostStmt ] .
3. RangeClause = [ ExpressionList “=” | IdentifierList “:=” ] “range” Expression .

**for 语句中临时变量是怎么回事？（为什么有时遍历赋值后，所有的值都等于最后一个元素）**
即使是短声明的变量，在 for 循环中也是`复用的`，这里的 v 一直 都是同一个临时变量，**所以&v 得到的地址一直都是相同的!!!**

```go
var a = make([]*int, 3)
for k, v := range []int{1, 2, 3} {
    a[k] = &v
}
for i := range a {
    fmt.Println(*a[i])
}
// result:
// 3
// 3
// 3
```

```go
// pushDown 函数
for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
    right := left + 1
    // ...
}
```

`for _,v := range` 循环切片时会拷贝切片里的元素，如果切片里保存的是值类型，例如很大的 struct，遍历的开销就会很大。
**这种时候需要使用 for i := range slice 的方式遍历，取值时使用 &slice[i]。**

## interface

1. []int 不能赋值给 []interface{}，可以赋值给 interface{}
   存放任意类型元素的切片可以使用 []interface{} 表示，但不能表示任意切片类型，即具体类型的切片无法转换为 []interface{} ，需要显示转换。
   https://blog.csdn.net/HaoDaWang/article/details/83931629

   ```go
   func foo(nums []interface{}) {}
   func bar(){
   	foo([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})  // error
   }
   ```

## 禁用 GC

```go
如果用的是指针写法，必要时禁止 GC，能加速不少
func init() { debug.SetGCPercent(-1) }
```

**https://github.dev/EndlessCheng/codeforces-go/blob/016834c19c4289ae5999988585474174224f47e2/copypasta**
Golang 卡常技巧（IO 之外的部分）
对于存在`海量小对象`的情况（如 trie, treap 等），使用 debug.SetGCPercent(-1) 来禁用 GC，能明显减少耗时
对于可以回收的情况（如 append 在超过 cap 时），使用 debug.SetGCPercent(-1) 虽然会减少些许耗时，但若有大量内存没被回收，会有 MLE 的风险
其他情况下使用 debug.SetGCPercent(-1) 对耗时和内存使用无明显影响
**对于多组数据的情况(力扣)，若禁用 GC 会 MLE，可在每组数据的开头或末尾调用 runtime.GC() 或 debug.FreeOSMemory() 手动 GC**
参考 https://draveness.me/golang/docs/part3-runtime/ch07-memory/golang-garbage-collector/
https://zhuanlan.zhihu.com/p/77943973
如果没有禁用 GC 但 MLE，可以尝试 1.19 新增的 `debug.SetMemoryLimit`
例如 debug.SetMemoryLimit(200<<20)，其中 200 可以根据题目的约束来修改

128MB ~1e7 个 int64
256MB ~3e7 个 int64
512MB ~6e7 个 int64
1GB ~1e8 个 int64

对于二维矩阵，以 make([][mx]int, n) 的方式使用，比 make([][]int, n) 嵌套 make([]int, m) 更高效（100MB 以上时可以快 ~150ms）
但需要注意这种方式可能会向 OS 额外申请一倍的内存
对比 https://codeforces.com/problemset/submission/375/118043978
https://codeforces.com/problemset/submission/375/118044262

**函数内的递归 lambda 会额外消耗非常多的内存**（~100MB / 1e6 递归深度）
写在 main 里面 + slice MLE https://codeforces.com/contest/767/submission/174193385
写在 main 里面 + array 257424KB https://codeforces.com/contest/767/submission/174194515
写在 main 外面 + slice 188364KB https://codeforces.com/contest/767/submission/174194380
写在 main 外面 + array 154500KB https://codeforces.com/contest/767/submission/174193693

## map 的 value 不可寻址 可使用指针类型代替

https://blog.csdn.net/qq_30505673/article/details/119919757
为什么？
map 会进行动静扩容，当进行扩大后，map 的 value 就会进行内存迁徙，
其地址发生变化，所以无法对这个 value 进行寻址。

## (结构体)内存对齐

如果存不下，下一个位置从 n%8==0 开始存储
每个字段都是 8 或者 4 的整数倍
建议在开发中，字段占用空间小的放在前面
