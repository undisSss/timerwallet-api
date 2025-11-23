import React, { useEffect, useState } from 'react'
import { getMe } from '../services/auth'
import { User } from '../types/user'

export default function Profile(){
  const [user, setUser] = useState<User | null>(null)

  useEffect(()=>{
    getMe().then(u => setUser(u)).catch(()=>{})
  },[])

  if(!user) return <div>Not logged in or loading...</div>

  return (
    <div>
      <h1>Profile</h1>
      <p><strong>ID:</strong> {user.id}</p>
      <p><strong>Name:</strong> {user.first_name} {user.last_name || ''}</p>
      <p><strong>Username:</strong> {user.username || '—'}</p>
      <p><strong>Roles:</strong> {(user.roles || []).join(', ') || 'user'}</p>
    </div>
  )
}
