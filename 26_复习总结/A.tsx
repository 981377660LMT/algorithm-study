import * as React from 'react'
import { useEffect, useState } from 'react'

interface AProps {}

const A: React.FC<AProps> = props => {
  const { children } = props
  const [count, setCount] = useState(0)

  useEffect(() => {
    document.title = `You clicked ${count} times`
  })

  return (
    <>
      <p>You clicked {count} times</p>
      <button onClick={() => setCount(count + 1)}>Click me</button>
    </>
  )
}

export { A }

type N = typeof Number
type n = number

type Key<T extends object, K = keyof T> = K extends any ? K : never

type OK = Key<Number>
