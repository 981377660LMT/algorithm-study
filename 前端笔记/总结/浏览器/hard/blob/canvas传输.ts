import axios from 'axios'

const canvas = document.createElement('canvas')
canvas.toBlob(blob => {
  if (blob) {
    const form = new FormData()
    form.append('canvasImage', blob, 'test.png')

    axios
      .post('url', form, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
      .then(res => console.log(res))
      .catch(err => console.log(err))
  }
})
