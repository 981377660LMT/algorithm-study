import * as React from 'react'
import { useEffect, useRef, useState } from 'react'
import * as ReactDom from 'react-dom'

interface IApp2Props {}

const App: React.FC<IApp2Props> = () => {
  const [count, setCount] = useState(0)
  const timerRef = useRef<number>()

  useEffect(() => {
    timerRef.current = window.setInterval(() => {
      setCount(prevData => prevData + 1)
    }, 1000)
    return () => {
      window.clearInterval(timerRef.current)
    }
  }, [])

  return <>{count}</>
}

const element = React.createElement(App)
ReactDom.render(element, document.querySelector('#app'))
