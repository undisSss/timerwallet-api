import React from 'react'
import { Link } from 'react-router-dom'
import { authenticateWithInitData } from '../services/auth'

export default function Home(){
  const onLogin = async () => {
    // In Telegram Mini Apps `window.Telegram.WebApp.initData` is provided.
    const initData = (window as any).Telegram?.WebApp?.initData || ''
    try{
      await authenticateWithInitData(initData)
      window.location.href = '/profile'
    }catch(err){
      alert('Auth failed')
    }
  }

  return (
    <div>
      <h1>TimerWallet Mini App</h1>
      <p>Welcome to the Mini App. Use the buttons below to navigate.</p>
      <div style={{display: 'flex', gap: '1rem'}}>
        <button onClick={onLogin}>Login with Telegram initData</button>
        <Link to="/game"><button>Go to Game</button></Link>
        <Link to="/profile"><button>Profile</button></Link>
      </div>
    </div>
  )
}
