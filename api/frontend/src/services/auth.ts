async function handleJsonResponse(res: Response){
  const text = await res.text()
  try{ return JSON.parse(text) }catch{return text}
}

export async function authenticateWithInitData(initData: string){
  const res = await fetch('/api/auth/telegram', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ initData })
  })
  if(!res.ok) throw new Error('Auth failed')
  return handleJsonResponse(res)
}

export async function refreshToken(){
  const res = await fetch('/api/auth/refresh', {
    method: 'POST',
    credentials: 'include'
  })
  if(!res.ok) throw new Error('Refresh failed')
  return handleJsonResponse(res)
}

export async function logout(){
  await fetch('/api/auth/logout', { method: 'POST', credentials: 'include' })
}

export async function getMe(){
  let res = await fetch('/api/me', { credentials: 'include' })
  if(res.status === 401){
    try{
      await refreshToken()
      res = await fetch('/api/me', { credentials: 'include' })
    }catch(e){
      throw new Error('Not authenticated')
    }
  }
  if(!res.ok) throw new Error('Not authenticated')
  return handleJsonResponse(res)
}
