import React, { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { userSummary } from '../api/client'
import { useAuth } from '../context/AuthContext'
import { clearTokens } from '../api/client'

export default function UserDashboard() {
  const nav = useNavigate()
  const { user, setUser } = useAuth()
  const [data, setData] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (user && user.role !== 'user') {
      nav('/403', { replace: true })
      return
    }
    async function load() {
      setLoading(true)
      setError('')
      try {
        const res = await userSummary()
        setData(res)
      } catch {
        setError('Failed to load user summary')
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [user, nav])

  function logout() {
    clearTokens()
    setUser(null)
    nav('/login', { replace: true })
  }

  return (
    <div className="page">
      <div className="card">
        <h1>User Dashboard</h1>
        <p className="muted">User-only API: /api/user/summary/</p>

        {loading ? <div className="center">Loading...</div> : null}
        {error ? <div className="error">{error}</div> : null}

        {data ? (
          <div className="kv">
            <div><span>Last Login</span><b>{data.last_login}</b></div>
            <div><span>Message</span><b>{data.message}</b></div>
          </div>
        ) : null}

        <div className="row">
          <button className="btn secondary" onClick={logout}>Logout</button>
        </div>
      </div>
    </div>
  )
}
