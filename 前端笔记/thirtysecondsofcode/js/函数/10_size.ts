size([1, 2, 3, 4, 5]) // 5
size('size') // 4
size({ one: 1, two: 2, three: 3 }) // 3

function size(val: unknown): number {
  if (Array.isArray(val)) return val.length
  else if (val != null && typeof val === 'object') return Object.keys(val).length
  else if (typeof val === 'string') return new Blob([val]).size
  return 0
}
