const useCount = () => {
  let count = 1

  const getCount = () => count

  const add = () => count++

  return {
    count, // 注意这里样获取count一个恒定的值(),因为字符串不可变
    countRef: { value: count }, // 注意这里样获取count一个恒定的值(),因为字符串不可变
    getCount, // 这样获取count才相当于获取属性
    add,
  }
}

const counter = useCount()
counter.add()
console.log(counter.count) // 1
console.log(counter.countRef.value) // 1
console.log(counter.getCount()) // 2

// 总结：只有函数才能动态访问另一个函数的内部变量
