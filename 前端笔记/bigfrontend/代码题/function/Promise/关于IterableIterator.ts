const iter = Array.from([() => {}, () => {}]).entries()

for (const [index, item] of iter) {
  console.log(index) // 0
  break
}

for (const [index, item] of iter) {
  console.log(index) // 1
  break
}

// 如果在数组共用iter(使用fill来共用iter的引用) 那么多次调用将耗尽iter 且每次调用不同
export {}
