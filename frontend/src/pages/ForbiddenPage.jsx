import React from 'react'
import { Link } from 'react-router-dom'

export default function ForbiddenPage() {
  return (
    <div className="page">
      <div className="card">
        <h1>403</h1>
        <p className="muted">You are not allowed to access this page.</p>
        <Link className="btn" to="/login">Back to Login</Link>
      </div>
    </div>
  )
}
