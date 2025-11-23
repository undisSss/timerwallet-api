import React, { createContext, useContext, useEffect, useState } from 'react'
import { getMe, authenticateWithInitData, refreshToken, logout as apiLogout } from '../services/auth'
import type { User } from '../types/user'

type AuthContextType = {
  user: User | null
  loading: boolean
  loginWithInitData: (initData: string) => Promise<void>
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider: React.FC<{children: React.ReactNode}> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  const loadMe = async () => {
    try{
      const u = await getMe()
      setUser(u)
    }catch(e){
      setUser(null)
    }finally{
      setLoading(false)
    }
  }

  useEffect(()=>{
    (async ()=>{
      try{
        await loadMe()
      }catch(_){}
      const iv = setInterval(async ()=>{
        try{ await refreshToken(); await loadMe() }catch(_){}
      }, 10 * 60 * 1000)
      return ()=> clearInterval(iv)
    })()
  }, [])

  const loginWithInitData = async (initData: string) => {
    setLoading(true)
    try{
      await authenticateWithInitData(initData)
      await loadMe()
    }finally{
      setLoading(false)
    }
  }

  const logout = async () => {
    try{ await apiLogout() }catch(_){}
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, loading, loginWithInitData, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(){
  const ctx = useContext(AuthContext)
  if(!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}

export default AuthContext
