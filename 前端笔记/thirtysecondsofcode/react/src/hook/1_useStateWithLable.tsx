import * as React from 'react'
import { useState, useDebugValue } from 'react'

interface UseStateWithLableProps {}

function useStateWithLabel<T>(initValue: T, label: string) {
  const [value, setValue] = useState(initValue)
  useDebugValue(`${label}:${value}`)
  return [value, setValue]
}

const UseStateWithLable: React.FC<UseStateWithLableProps> = props => {
  const [value] = useStateWithLabel(0, 'counteer')

  return (
    <>
      <p>{value}</p>
    </>
  )
}

export { UseStateWithLable }
