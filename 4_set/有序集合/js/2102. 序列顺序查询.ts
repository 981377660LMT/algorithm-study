import { TreapMultiSet } from './Treap'

interface Spot {
  name: string
  score: number
}

const leftInf: Spot = { name: 'a', score: Infinity }
const rightInf: Spot = { name: 'z', score: -Infinity }
const compareFn = (a: Spot, b: Spot) => -(a.score - b.score) || a.name.localeCompare(b.name)

class SORTracker {
  private readonly sortedList: TreapMultiSet<Spot> = new TreapMultiSet(compareFn, leftInf, rightInf)
  private id = 0

  add(name: string, score: number): void {
    this.sortedList.add({ name, score })
  }

  get(): string {
    return this.sortedList.at(this.id++)!.name
  }
}

const st = new SORTracker()
st.add('bradford', 2)
st.add('bradford', 3)
st.add('bradford', 3)
st.add('best', 100)
st.add('best2', 100)

console.log(st.get())
console.log(st.get())
console.log(st.get())
console.log(st.get())
console.log(st.get())
console.log(st.get())
