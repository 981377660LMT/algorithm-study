const newPost = {
  userId: 1,
  id: 1337,
  title: 'Foo',
  body: 'bar bar bar',
}
const data = JSON.stringify(newPost)
httpPost('https://jsonplaceholder.typicode.com/posts', data, console.log) /*
Logs: {
  "userId": 1,
  "id": 1337,
  "title": "Foo",
  "body": "bar bar bar"
}
*/

function httpPost(
  url: string,
  data: string,
  onload: (...args: any[]) => void,
  onerror: (...args: any[]) => void = console.error
) {
  const req = new XMLHttpRequest()
  req.open('POST', url, true)
  // setRequestHeader 方法设置 HTTP 请求头的值
  req.setRequestHeader('Content-type', 'application/json; charset=utf-8')
  req.onload = () => onload(req.responseText)
  req.onerror = () => onerror(req)
  req.send(data)
}
