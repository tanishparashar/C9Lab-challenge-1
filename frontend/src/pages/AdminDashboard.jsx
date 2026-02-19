import React, { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { adminStats } from '../api/client'
import { useAuth } from '../context/AuthContext'
import { clearTokens } from '../api/client'

export default function AdminDashboard() {
  const nav = useNavigate()
  const { user, setUser } = useAuth()
  const [data, setData] = useState(null)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (user && user.role !== 'admin') {
      nav('/403', { replace: true })
      return
    }
    async function load() {
      setLoading(true)
      setError('')
      try {
        const res = await adminStats()
        setData(res)
      } catch {
        setError('Failed to load admin stats')
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
        <h1>Admin Dashboard</h1>
        <p className="muted">Admin-only API: /api/admin/stats/</p>

        {loading ? <div className="center">Loading...</div> : null}
        {error ? <div className="error">{error}</div> : null}

        {data ? (
          <div className="kv">
            <div><span>Total Users</span><b>{data.total_users}</b></div>
            <div><span>Active Users</span><b>{data.active_users}</b></div>
          </div>
        ) : null}

        <div className="row">
          <button className="btn secondary" onClick={logout}>Logout</button>
        </div>
      </div>
    </div>
  )
}
