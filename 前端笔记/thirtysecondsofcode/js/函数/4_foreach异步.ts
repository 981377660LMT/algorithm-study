const asyncUppercase = (item: string) =>
  new Promise(resolve =>
    setTimeout(() => resolve(item.toUpperCase()), Math.floor(Math.random() * 1000))
  )

const uppercaseItems = async () => {
  const items = ['a', 'b', 'c']
  // "await" 对此表达式的类型没有影响。
  await items.forEach(async item => {
    const uppercaseItem = await asyncUppercase(item)
    console.log(uppercaseItem)
  })

  console.log('Items processed')
}

uppercaseItems()
