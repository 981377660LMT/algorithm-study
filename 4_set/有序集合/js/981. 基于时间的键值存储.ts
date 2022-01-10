import { TreeSet } from './TreeSet'

interface Info {
  value: string
  timestamp: number
}

class TimeMap {
  private record: Map<string, TreeSet<Info>>

  constructor() {
    this.record = new Map()
  }

  set(key: string, value: string, timestamp: number): void {
    !this.record.has(key) &&
      this.record.set(key, new TreeSet<Info>([], (a: Info, b: Info) => a.timestamp - b.timestamp))
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
