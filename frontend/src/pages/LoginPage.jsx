import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { login as apiLogin, me as apiMe } from '../api/client'
import { useAuth } from '../context/AuthContext'

function isValidEmail(v) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v)
}

export default function LoginPage() {
  const nav = useNavigate()
  const { setUser } = useAuth()

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [rememberMe, setRememberMe] = useState(true)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function onSubmit(e) {
    e.preventDefault()
    setError('')

    const em = email.trim()
    if (!isValidEmail(em)) {
      setError('Please enter a valid email')
      return
    }
    if (!password) {
      setError('Password is required')
      return
    }

    setLoading(true)
    try {
      await apiLogin(em, password, rememberMe)
      const u = await apiMe()
      setUser(u)

      if (u.role === 'admin') nav('/admin-dashboard', { replace: true })
      else nav('/user-dashboard', { replace: true })
    } catch (err) {
      const msg = err?.response?.data?.detail || err?.response?.data?.error || 'Login failed'
      setError(msg)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="page">
      <div className="card">
        <h1>Login</h1>
        <p className="muted">Project 2: Role-based auth + refresh token</p>

        <form onSubmit={onSubmit} className="form">
          <label>
            Email
            <input value={email} onChange={(e) => setEmail(e.target.value)} type="email" placeholder="admin@company.com" required />
          </label>

          <label>
            Password
            <input value={password} onChange={(e) => setPassword(e.target.value)} type="password" placeholder="••••••••" required />
          </label>

          <label className="checkbox">
            <input type="checkbox" checked={rememberMe} onChange={(e) => setRememberMe(e.target.checked)} />
            Remember me
          </label>

          {error ? <div className="error">{error}</div> : null}

          <button className="btn" disabled={loading}>
            {loading ? 'Signing in...' : 'Sign In'}
          </button>

          <div className="hint">
            Test users:
            <ul>
              <li>admin@example.com / admin123</li>
              <li>user@example.com / user123</li>
            </ul>
          </div>
        </form>
      </div>
    </div>
  )
}
