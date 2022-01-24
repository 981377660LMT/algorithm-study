const download = (fileName: string, blob: Blob) => {
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = fileName
  link.click()
  link.remove()
  URL.revokeObjectURL(link.href)
}

const downloadBtn = document.querySelector('#downloadBtn')!
downloadBtn.addEventListener('click', () => {
  const fileName = 'blob.txt'
  const myBlob = new Blob(['一文彻底掌握 Blob Web API'], { type: 'text/plain' })
  download(fileName, myBlob)
})
