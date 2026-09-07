import axios from 'axios'
import { useAuthStore } from '../stores/auth'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api',
  timeout: 10000,
  // withCredentials: true
})

api.interceptors.request.use((config) => {
  const auth = useAuthStore()
  
  // ✅ បន្ថែម x-api-key
  config.headers['x-api-key'] = import.meta.env.API_KEY_SECRET || 'X7mK9qR2vLp8Nf4TzY6cHd3WsAj5BuEeG1rQn8MxKp7Vt2CyL9sDf4JhUwZk3NaRb'
  config.headers['X-Admin-Token'] = import.meta.env.ADMIN_SECRET || 'X7mK9qR2vLp8Nf4TzY6cHd3WsAj5BuEeG1rQn8MxKp7Vt2CyL9sDf4JhUwZk3NaRb៣៩៨៩៩ឪឧ' 
  // JWT Token (ប្រសិនបើមាន)
  if (auth.accessToken) {
    config.headers.Authorization = `Bearer ${auth.accessToken}`
  }

  return config
})

let isRefreshing = false
let failedQueue = []

function processQueue(error, token = null) {
  failedQueue.forEach(({ resolve, reject }) => {
    if (error) {
      reject(error)
    } else {
      resolve(token)
    }
  })
  failedQueue = []
}


api.interceptors.response.use(
  (res) => res,
  async (error) => {
    const auth = useAuthStore()
    const original = error.config

    if (error.response?.status === 401 && !original._retry) {
      if (isRefreshing) {
        // queue this request until the in-flight refresh finishes
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        }).then((token) => {
          original.headers.Authorization = `Bearer ${token}`
          return api(original)
        }).catch((err) => Promise.reject(err))
      }

      original._retry = true
      isRefreshing = true

      try {
        await auth.refresh()
        const newToken = auth.accessToken
        processQueue(null, newToken)
        original.headers.Authorization = `Bearer ${newToken}`
        return api(original)
      } catch (refreshError) {
        processQueue(refreshError, null)
        auth.logout()
        window.location.href = '/login'
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }

    return Promise.reject(error)
  }
)


export default api