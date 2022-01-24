const file = new File(['a'.repeat(1000000)], 'test.txt')

const chunkSize = 40000
const url = 'https://httpbin.org/post'

async function chunkedUpload() {
  for (let start = 0; start < file.size; start += chunkSize) {
    const chunk = file.slice(start, start + chunkSize)
    const fd = new FormData()
    fd.append('data', chunk)

    await fetch(url, { method: 'post', body: fd }).then(res => res.text())
  }
}

export {}
