import { useRequest } from 'ahooks'

function App() {
  const { data, error, loading } = useRequest(() => fetch('https://api.github.com').then(res => res.json()))

  if (loading) return <div>Loading...</div>
  if (error) return <div>Error!</div>
  return <pre>{JSON.stringify(data, null, 2)}</pre>
}

export default App
