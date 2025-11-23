import React from 'react'

export default function Game(){
  const startGame = () => {
    alert('Game started (placeholder)')
  }

  return (
    <div>
      <h1>Game</h1>
      <p>This is a placeholder game page. Integrate your game here.</p>
      <button onClick={startGame}>Start</button>
    </div>
  )
}
