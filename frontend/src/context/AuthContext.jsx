import React, { createContext, useContext, useEffect, useMemo, useState } from 'react'
import { clearTokens, getAccessToken, me as apiMe } from '../api/client'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)

  const isAuthenticated = !!getAccessToken()

  useEffect(() => {
    async function load() {
      try {
        if (getAccessToken()) {
          const u = await apiMe()
          setUser(u)
        }
      } catch {
        clearTokens()
        setUser(null)
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [])

  const value = useMemo(() => ({ user, setUser, loading, isAuthenticated }), [user, loading, isAuthenticated])
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  return useContext(AuthContext)
}
