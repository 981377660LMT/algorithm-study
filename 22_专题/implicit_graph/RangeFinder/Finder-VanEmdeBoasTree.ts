import { VanEmdeBoasTree } from '../../../24_高级数据结构/珂朵莉树/VanEmdeBoasTree'

if (require.main === module) {
  const van = new VanEmdeBoasTree()
  console.log(van.min, van.max, van.size)
  van.insert(3)
  van.insert(1)
  van.insert(2)
  console.log(van.has(1))
  console.log(van.has(2))
  console.log(van.has(3))
  console.log(van.toString())

  const n = 2e5
  const set2 = new VanEmdeBoasTree()
  console.time('VanEmdeBoasTree')
  for (let i = 0; i < n; i++) {
    set2.insert(i)
    set2.next(i)
    set2.prev(i)
    set2.has(i)
    set2.erase(i)
    set2.insert(i)
  }
  console.timeEnd('VanEmdeBoasTree') // 360ms
}
