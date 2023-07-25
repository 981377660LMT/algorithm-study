import { SortedList } from './_SortedList'

class SORTracker {
  private readonly _sl = new SortedList<[score: number, name: string]>(
    (a, b) => -(a[0] - b[0]) || a[1].localeCompare(b[1])
  )

  private _ptr = 0

  add(name: string, score: number): void {
    this._sl.add([score, name])
  }

  get(): string {
    return this._sl.at(this._ptr++)![1]
  }
}

/**
 * Your SORTracker object will be instantiated and called as such:
 * var obj = new SORTracker()
 * obj.add(name,score)
 * var param_2 = obj.get()
 */
