import { TreapMultiSet as SortedList } from '../Treap'

describe('treap', () => {
  describe('TreapMultiSet', () => {
    it('should dedupe', () => {
      const sortedList = new SortedList().add(0, 0, 1, 1, 2)
      expect([...sortedList.values()]).toEqual([0, 0, 1, 1, 2])
    })

    it('should support add/has/delete', () => {
      const sortedList = new SortedList().add(0, 1, 2, 3, 4, 5)
      for (let i = 0; i <= 5; i++) expect(sortedList.has(i)).toEqual(true)
      expect(sortedList.has(10)).toEqual(false)
      expect(sortedList.has(-10)).toEqual(false)
      sortedList.delete(5)
      expect(sortedList.has(5)).toEqual(false)
    })

    it('should support bisectLeft/bisectRight', () => {
      const sortedList = new SortedList().add(2, 3, 3)
      expect(sortedList.bisectLeft(1.5)).toEqual(0)
      expect(sortedList.bisectLeft(2)).toEqual(0)
      expect(sortedList.bisectLeft(2.5)).toEqual(1)
      expect(sortedList.bisectLeft(3)).toEqual(1)
      expect(sortedList.bisectLeft(3.5)).toEqual(3)
      expect(sortedList.bisectRight(1.5)).toEqual(0)
      expect(sortedList.bisectRight(2)).toEqual(1)
      expect(sortedList.bisectRight(2.5)).toEqual(1)
      expect(sortedList.bisectRight(3)).toEqual(3)
      expect(sortedList.bisectRight(3.5)).toEqual(3)
    })

    it('should support indexOf/lastIndexOf', () => {
      const sortedList = new SortedList().add(1, 2, 3, 3, 3, 4)
      expect(sortedList.indexOf(1)).toEqual(0)
      expect(sortedList.indexOf(3)).toEqual(2)
      expect(sortedList.indexOf(5)).toEqual(-1)
      expect(sortedList.lastIndexOf(1)).toEqual(0)
      expect(sortedList.lastIndexOf(3)).toEqual(4)
      expect(sortedList.lastIndexOf(5)).toEqual(-1)
    })

    it('should support at/first/last', () => {
      const sortedList = new SortedList().add(1, 2, 3, 3, 3)
      expect(sortedList.at(0)).toEqual(1)
      expect(sortedList.at(1)).toEqual(2)
      expect(sortedList.at(10)).toEqual(undefined)
      expect(sortedList.at(-1)).toEqual(3)
      expect(sortedList.at(-4)).toEqual(2)
      expect(sortedList.at(-10)).toEqual(undefined)
      expect(sortedList.first()).toEqual(1)
      expect(sortedList.last()).toEqual(3)
    })

    it('should support ceil/floor/higher/lower', () => {
      const sortedList = new SortedList().add(1, 2, 3, 4, 5)
      expect(sortedList.ceil(1)).toEqual(1)
      expect(sortedList.ceil(5.5)).toEqual(undefined)
      expect(sortedList.floor(3.5)).toEqual(3)
      expect(sortedList.floor(1)).toEqual(1)
      expect(sortedList.higher(1)).toEqual(2)
      expect(sortedList.higher(5)).toEqual(undefined)
      expect(sortedList.lower(1)).toEqual(undefined)
      expect(sortedList.lower(4.5)).toEqual(4)
    })

    it('should support shift/pop', () => {
      const sortedList = new SortedList().add(1, 2, 3, 4, 5)
      expect(sortedList.shift()).toEqual(1)
      expect(sortedList.size).toEqual(4)
      expect(sortedList.pop()).toEqual(5)
      expect(sortedList.size).toEqual(3)
      expect(sortedList.pop(1)).toEqual(3)
      expect(sortedList.size).toEqual(2)
    })

    it('should support count', () => {
      const sortedList = new SortedList().add(0, 1, 1, 2)
      expect(sortedList.count(-1)).toEqual(0)
      expect(sortedList.count(0)).toEqual(1)
      expect(sortedList.count(1)).toEqual(2)
      expect(sortedList.count(2)).toEqual(1)
    })

    it('should sort the values correctly', () => {
      for (let i = 0; i < 20; i++) {
        const arr = [...Array(100)].map(() => Math.random())
        const sortedList = new SortedList().add(...arr)
        expect([...sortedList.keys()]).toEqual(arr.sort((a, b) => a - b))
        expect([...sortedList.values()]).toEqual(arr.sort((a, b) => a - b))
        expect([...sortedList.rvalues()]).toEqual(arr.sort((a, b) => b - a))
      }

      const sortedList = new SortedList().add(1, 2, 3, 3, 3)
      expect([...sortedList.entries()]).toEqual([...sortedList.values()].map((v, i) => [i, v]))
    })

    it('should support custom compare function', () => {
      const bigHead = new SortedList((a, b) => -(a - b), Infinity, -Infinity).add(1, 2, 3)
      expect([...bigHead.values()]).toEqual([3, 2, 1])

      interface Student {
        name: string
        score: number
      }

      const leftInf: Student = { name: 'a', score: Infinity }
      const rightInf: Student = { name: 'z', score: -Infinity }
      const compareFn = (a: Student, b: Student): number =>
        -(a.score - b.score) || a.name.localeCompare(b.name)
      const testRank = new SortedList<Student>(compareFn, leftInf, rightInf)

      const person1 = { name: 'Bob', score: 2 }
      const person2 = { name: 'Alice', score: 2 }
      const person3 = { name: 'Foo', score: 3 }
      const person4 = { name: 'Bar', score: 0 }

      testRank.add(person1).add(person2).add(person3).add(person4)
      expect([...testRank.values()]).toEqual([person3, person2, person1, person4])
    })
  })
})
