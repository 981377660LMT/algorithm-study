import { TreapMultiSet } from './Treap'

interface Info {
  value: string
  timestamp: number
}

class TimeMap {
  private record: Map<string, TreapMultiSet<Info>>

  constructor() {
    this.record = new Map()
  }

  set(key: string, value: string, timestamp: number): void {
    !this.record.has(key) &&
      this.record.set(
        key,
        new TreapMultiSet<Info>(
          (a: Info, b: Info) => a.timestamp - b.timestamp,
          {
            value: '',
            timestamp: -Infinity,
          },
          {
            value: '',
            timestamp: Infinity,
          }
        )
      )
    this.record.get(key)!.add({ value, timestamp })
  }

  get(key: string, timestamp: number): string {
    if (!this.record.has(key)) return ''
    const treeSet = this.record.get(key)!
    const cand = treeSet.floor({
      timestamp,
      value: '',
    })

    return cand?.value ?? ''
  }
}

const tm = new TimeMap()
tm.set('foo', 'bar', 1)
console.log(tm.get('foo', 1))
