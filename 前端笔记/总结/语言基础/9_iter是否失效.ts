const map = new Map<string, string>([])

map.set('name', 'cmnx').set('sex', 'm').set('age', '666').set('weight', '777')

const iter = map.entries()

map.delete('sex')

console.log([...iter])
