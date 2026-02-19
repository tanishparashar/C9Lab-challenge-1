import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

export const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 15000,
})

function getStorage(rememberMe) {
  return rememberMe ? localStorage : sessionStorage
}

export function getAccessToken() {
  return localStorage.getItem('access_token') || sessionStorage.getItem('access_token')
}

export function getRefreshToken() {
  return localStorage.getItem('refresh_token') || sessionStorage.getItem('refresh_token')
}

export function clearTokens() {
  localStorage.removeItem('access_token')
  localStorage.removeItem('refresh_token')
  sessionStorage.removeItem('access_token')
  sessionStorage.removeItem('refresh_token')
}

export function setTokens({ access, refresh, rememberMe }) {
  clearTokens()
  const storage = getStorage(rememberMe)
  storage.setItem('access_token', access)
  storage.setItem('refresh_token', refresh)
}

let isRefreshing = false
let pendingQueue = []

function processQueue(error, token = null) {
  pendingQueue.forEach(p => (error ? p.reject(error) : p.resolve(token)))
  pendingQueue = []
}

api.interceptors.request.use((config) => {
  const token = getAccessToken()
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const original = error.config
    const status = error?.response?.status

    if (status === 401 && !original?._retry) {
      original._retry = true
      const refresh = getRefreshToken()
      if (!refresh) {
        clearTokens()
        return Promise.reject(error)
      }

      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          pendingQueue.push({ resolve, reject })
        }).then((newToken) => {
          original.headers.Authorization = `Bearer ${newToken}`
          return api(original)
        })
      }

      isRefreshing = true
      try {
        const res = await api.post('/api/auth/refresh/', { refresh })
        const newAccess = res.data.access
        // preserve where tokens were stored
        if (localStorage.getItem('refresh_token')) localStorage.setItem('access_token', newAccess)
        if (sessionStorage.getItem('refresh_token')) sessionStorage.setItem('access_token', newAccess)

        processQueue(null, newAccess)
        original.headers.Authorization = `Bearer ${newAccess}`
        return api(original)
      } catch (e) {
        processQueue(e, null)
        clearTokens()
        return Promise.reject(e)
      } finally {
        isRefreshing = false
      }
    }

    return Promise.reject(error)
  }
)

export async function login(email, password, rememberMe) {
  const res = await api.post('/api/auth/login/', { email, password })
  setTokens({ access: res.data.access, refresh: res.data.refresh, rememberMe })
  return res.data
}

export async function me() {
  const res = await api.get('/api/auth/me/')
  return res.data
}

export async function adminStats() {
  const res = await api.get('/api/admin/stats/')
  return res.data
}

export async function userSummary() {
  const res = await api.get('/api/user/summary/')
  return res.data
}
