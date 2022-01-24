function readAsArrayBuffer(file: Blob, start = 0, end = 2) {
  return new Promise<ArrayBuffer>((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => {
      resolve(reader.result as ArrayBuffer)
    }
    reader.onerror = reject
    reader.readAsArrayBuffer(file.slice(start, end))
  })
}

// 对于 PNG 类型的图片来说，该文件的前 8 个字节是 0x89 50 4E 47 0D 0A 1A 0A。
// 因此，我们在检测已选择的文件是否为 PNG 类型的图片时，
// 只需要读取前 8 个字节的数据，并逐一判断每个字节的内容是否一致。
function check(headers: number[]) {
  return (buffers: ArrayBuffer, options = { offset: 0 }) =>
    headers.every((header, index) => header === buffers[options.offset + index])
}

export {}

const isPNG = check([0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a]) // PNG图片对应的魔数
const realFileElement = document.querySelector('#realFileType') as HTMLInputElement

async function handleChange(event: Event) {
  const file = (event.target as HTMLInputElement).files![0]
  const buffer = await readAsArrayBuffer(file, 0, 8)
  const uint8Array = new Uint8Array(buffer)
  realFileElement.innerText = `${file.name}文件的类型是：${
    isPNG(uint8Array) ? 'image/png' : file.type
  }`
}
