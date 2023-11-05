import { SegmentTreeDivideAndConquerCopy } from '../SegmentTreeDivideAndConquerCopy'
import { SegmentTreeDivideAndConquerUndo } from '../SegmentTreeDivideAndConquerUndo'
import { SweepLine } from '../SweepLine'
import { mutateWithoutOneCopy } from '../mutateWithOutOneCopy'

function test1(): void {
  let mutate = 0
  let undo = 0
  let query = 0

  const dc = new SegmentTreeDivideAndConquerUndo({
    mutate(id) {
      mutate++
    },
    undo() {
      undo++
    },
    query(id) {
      query++
    }
  })

  const n = 1e5
  for (let i = 0; i < n; i++) {
    dc.addMutation(0, i)
    dc.addMutation(i, n)
  }
  for (let i = 0; i < n; i++) {
    dc.addQuery(i)
  }

  dc.run()
  console.log(mutate, undo, query) // 2056006 2056006 100000
}

test1()

function test2(): void {
  let mutate = 0
  let copy = 0
  let query = 0

  const dc = new SegmentTreeDivideAndConquerCopy(
    { value: 1 },
    {
      mutate(id) {
        mutate++
      },
      copy() {
        copy++
        return { value: 1 }
      },
      query(id) {
        query++
      }
    }
  )

  const n = 1e5
  for (let i = 0; i < n; i++) {
    dc.addMutation(0, i)
    dc.addMutation(i, n)
  }
  for (let i = 0; i < n; i++) {
    dc.addQuery(i)
  }

  dc.run()
  console.log(mutate, copy, query) // 2056006 231070 100000
}

test2()

function test3(): void {
  let mutate = 0
  let copy = 0
  let query = 0

  const n = 1e5

  mutateWithoutOneCopy({ value: 1 }, 0, n, {
    mutate(state, index) {
      mutate++
    },
    copy() {
      copy++
      return { value: 1 }
    },
    query(state, index) {
      query++
    }
  })

  console.log(mutate, copy, query) // 1668928 199998 100000
}

test3()

function test4(): void {
  let mutate = 0
  let remove = 0
  let query = 0

  const n = 1e5

  const sweepLine = new SweepLine({
    mutate(mutationId) {
      mutate++
    },
    remove(mutationId) {
      remove++
    },
    query(mutationId) {
      query++
    }
  })

  for (let i = 0; i < n; i++) {
    sweepLine.addMutation(i, i + 1, i)
  }
  for (let i = 0; i < 2 * n; i++) {
    sweepLine.addQuery(i, i)
  }
  sweepLine.run()

  console.log(mutate, remove, query) // 100000 100000 200000
}

test4()
