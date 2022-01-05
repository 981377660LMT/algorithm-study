const iter = Array.from([() => {}, () => {}]).entries()

for (const [index, _item] of iter) {
  console.log(index) // 0
  break
}

for (const [index, _item] of iter) {
  console.log(index) // 1
  break
}

// 如果在数组共用iter(使用fill来共用iter的引用) 那么多次调用将耗尽iter 且每次调用不同
export {}

// iterator 是会被消耗的 而 iterable 不会
