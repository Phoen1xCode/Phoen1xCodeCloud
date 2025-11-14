import axios from 'axios'

const api = axios.create({
  baseURL: '/api'
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export const authAPI = {
  register: (data) => api.post('/register', data),
  login: (data) => api.post('/login', data)
}

export const shareAPI = {
  uploadFile: (formData) => api.post('/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  }),
  createText: (data) => api.post('/text', data),
  getShare: (code) => api.get(`/share/${code}`),
  listShares: () => api.get('/shares'),
  deleteShare: (code) => api.delete(`/share/${code}`)
}

export const adminAPI = {
  getStats: () => api.get('/admin/stats'),
  listAllShares: () => api.get('/admin/shares'),
  listUsers: () => api.get('/admin/users')
}

export default api
