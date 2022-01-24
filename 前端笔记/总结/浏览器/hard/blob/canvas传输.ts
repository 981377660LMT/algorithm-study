import axios from 'axios'

const canvas = document.createElement('canvas')
const dataURL = canvas.toDataURL()
const form = new FormData()
const myimg = new Blob([dataURL])
form.append('', myimg)

axios
  .post('url', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  .then(res => console.log(res))
  .catch(err => console.log(err))
