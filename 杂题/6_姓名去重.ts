const nameDeduplicate = (names: string[]): string[] => {
  if (names.length <= 1) return names
  return [...new Set(names.map(name => name.toLowerCase()))]
}

console.log(
  nameDeduplicate([
    'James',
    'james',
    'Bill Gates',
    'bill Gates',
    'Hello World',
    'HELLO WORLD',
    'Helloworld',
  ])
)

export {}
