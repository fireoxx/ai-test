import axios from 'axios'

const service = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 10000
})

service.interceptors.response.use(
  (response) => response.data,
  (error) => {
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

export default service
