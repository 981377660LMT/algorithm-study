import * as React from 'react'
import * as ReactDom from 'react-dom'

export interface IAppProps {}

export interface IAppState {
  val: number
}

class App1 extends React.Component<IAppProps, IAppState> {
  constructor(props: IAppProps) {
    super(props)
    this.state = {
      val: 0,
    }
  }

  override componentDidMount() {
    this.setState({ val: this.state.val + 1 }, () => console.log(this.state.val, '1111')) // 1
    console.log(this.state.val) // 0

    this.setState({ val: this.state.val + 1 }, () => console.log(this.state.val, '222')) // 1
    console.log(this.state.val) // 0

    setTimeout(() => {
      this.setState({ val: this.state.val + 1 })
      console.log(this.state.val) // 2

      this.setState({ val: this.state.val + 1 })
      console.log(this.state.val) // 3
    }, 0)
  }

  override render() {
    return <div>{this.state.val}</div> // 3
  }
}
//////////////////////////////////////////////////////////////////////
interface IApp2Props {}

const App2: React.FC<IApp2Props> = props => {
  const [val, setVal] = React.useState(0)

  React.useEffect(() => {
    setVal(val => val + 1)
    console.log(val) // 闭包 输出0
    setVal(val => val + 1)
    console.log(val) // 闭包 输出0
    setTimeout(() => {
      setVal(val => val + 1)
      console.log(val) // 闭包 输出0
      setVal(val => val + 1)
      console.log(val) // 闭包 输出0
    }, 0)
  }, [])

  return <div>{val}</div> // 4
}

const element = React.createElement(App2)
ReactDom.render(element, document.querySelector('#app'))

// 类组件：
//0 0 2 3
// 出于性能考虑，React 可能会把多个 setState() 调用合并成一个调用。
// 1、第一次和第二次都是在 react 自身生命周期内，触发时 isBatchingUpdates 为 true，所以并不会直接执行更新 state，而是加入了 dirtyComponents，所以打印时获取的都是更新前的状态 0。
// 2、两次 setState 时，获取到 this.state.val 都是 0，所以执行时都是将 0 设置成 1，在 react 内部会被合并掉，只执行一次。设置完成后 state.val 值为 1。
// 如果不想被合并呢?
// 要解决这个问题，可以让 setState() 接收一个函数而不是一个对象
// state: ((prevState: Readonly<S>, props: Readonly<P>) => (Pick<S, K> | S | null))
// 3、setTimeout 中的代码，触发时 isBatchingUpdates 为 false，所以能够直接进行更新，所以连着输出 2，3。
// 输出： 0 0 2 3

// 函数组件
// 0 0 0 0
