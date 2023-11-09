import { parallelBinarySearch } from '../parallelBinarySearch'
import { parallelBinarySearchUndo } from '../parallelBinarySearchUndo'

function testCount1(): void {
  const n = 1e2
  const q = 1e5
  let mutateCount = 0
  let resetCount = 0
  let predicateCount = 0

  parallelBinarySearch(n, q, {
    mutate() {
      mutateCount++
    },
    reset() {
      resetCount++
    },
    predicate() {
      predicateCount++
      return false
    }
  })

  console.log({ mutateCount, resetCount, predicateCount })
  // { mutateCount: 603, resetCount: 7, predicateCount: 800000 }
}

function testCount2(): void {
  const n = 1e2
  const q = 1e5
  let mutateCount = 0
  let undoCount = 0
  let predicateCount = 0

  parallelBinarySearchUndo(n, q, {
    mutate() {
      mutateCount++
    },
    undo() {
      undoCount++
    },
    predicate() {
      predicateCount++
      return false
    }
  })
  console.log({ mutateCount, undoCount, predicateCount })
  // { mutateCount: 603, undoCount: 603, predicateCount: 800000 }
}

testCount1()
testCount2()
