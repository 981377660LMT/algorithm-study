const file = document.querySelector('#file') as HTMLInputElement
const img = document.querySelector('#img') as HTMLImageElement

// // 方式1：URL.createObjectURL()
// file.onchange = () => {
//   const data = URL.createObjectURL(file.files![0])
//   // console.log(data)
//   // blob:http://localhost:1234/12a9d830-a83c-4f23-b5a3-cdd6d344f05a

//   img.src = data

//   // 此方式还需要在安全的时机主动释放，否则占用内存
//   // URL.revokeObjectURL(data)
// }

// 方式2：FileReader
file.onchange = () => {
  // 不同于 nodejs 的 fs 可以读取磁盘文件，浏览器中的 FileReader 只能操作用户选择的文件
  const reader = new FileReader()
  reader.readAsDataURL(file.files![0])
  reader.onload = () => {
    img.src = reader.result as string
    console.log(reader.result)
    // data:image/jpeg;base64,/9j/2wCEAAgGBgcGBQg...
  }
}
