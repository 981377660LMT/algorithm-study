interface ICount {
  (): number
  reset: () => void
}

function useCount(): ICount {
  let state = 0

  const count: ICount = () => {
    state++
    return state
  }

  count.reset = () => {
    state = 0
  }

  return count
}

const count = useCount()
console.log(count()) // 1
console.log(count()) // 2
console.log(count()) // 3
count.reset()
console.log(count()) // 1
console.log(count()) // 2
console.log(count()) // 3

export {}
