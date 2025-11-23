import React from 'react'
import { Link } from 'react-router-dom'

export default function Header(){
  return (
    <header style={{padding: '0.5rem 1rem', borderBottom: '1px solid #ddd'}}>
      <nav style={{display: 'flex', gap: '1rem', alignItems: 'center'}}>
        <Link to="/">Home</Link>
        <Link to="/game">Game</Link>
        <Link to="/profile">Profile</Link>
      </nav>
    </header>
  )
}
