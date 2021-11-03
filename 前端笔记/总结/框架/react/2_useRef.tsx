import * as React from 'react'
import { useEffect, useRef } from 'react'
import * as ReactDom from 'react-dom'

interface IApp2Props {}

const App: React.FC<IApp2Props> = () => {
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    console.log(inputRef.current)
    inputRef.current?.focus()
    return () => {}
  }, [])

  return (
    <>
      <Input ref={inputRef}>点我</Input>
    </>
  )
}

interface IButtonProps {
  children: React.ReactNode
}

const Input = React.forwardRef<HTMLInputElement, IButtonProps>((props, ref) => {
  const { children } = props
  return (
    <div>
      {children}
      <input ref={ref} />
    </div>
  )
})

const element = React.createElement(App)
ReactDom.render(element, document.querySelector('#app'))
