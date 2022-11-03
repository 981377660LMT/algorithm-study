https://qiita.com/ktateish/items/ab2df3e0864d2e931bf2

- AtCoder のジャッジシステムは 64bit 環境であり int 型が 64bit サイズになっているため、int64 型をわざわざ使い分けなくて良いのは楽です。
- slice の先頭と末尾からであれば要素を O(1) で削除できます

## golang 切片

1. 切片的[:]操作取的是一段 subarray，会共享原切片底层数组(这一点与 python 不同)
2. append 进行扩容时，什么时候会共享原切片底层数组？什么时候会使用新数组来创建切片？
   append 之后元素个数**小于等于原切片剩余容量时，会共享原切片底层数组，否则会扩容使用新数组来创建切片**
3. 如何对切片进行浅拷贝?

   > 浅拷贝指的是，把变量里面存的内存地址拷贝了，所指向的真实值并没拷贝。

   第一种方法是 append+解构 (类似 js 的 [...nums] )

   ```go
   func main() {
   	slice1 := []int{1, 2, 3, 4, 5}
   	sliceCopy := append([]int{}, slice1...)
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

4. 切片内存储类似 tuple 的对时，可以用结构体实现
5. 切片初始化
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

6. make 和 new 的区别
   make 用于 chan map slice 的初始化，new 用于类型的属性初始化
   他们都返回的是指针(`注意 chan/map/slice 本身就是指针`)
7. make 初始化 struct 是`深层次的初始化`

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
2. 闭包 + 返回接口
