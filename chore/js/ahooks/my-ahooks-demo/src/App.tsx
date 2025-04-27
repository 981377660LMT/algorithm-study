import { useRequest } from 'ahooks'
import Mock from 'mockjs'

function getUsername(): Promise<string> {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve(Mock.mock('@name'))
    }, 1000)
  })
}

function App() {
  const { data, error, loading, mutate } = useRequest(getUsername)

  if (error) {
    return <div>failed to load</div>
  }
  if (loading) {
    return <div>loading...</div>
  }
  return <div>Username: {data}</div>
}

export default App
