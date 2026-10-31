// Holds "who is logged in" for the whole app.
import { createContext, useCallback, useContext, useEffect, useState, type ReactNode } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { api, setUnauthorizedHandler, tokenStore } from '../lib/api'
import type { LoginResponse, User } from '../lib/types'

interface AuthState {
  user: User | null
  isLoading: boolean // true while we check a saved token on page load
  login: (email: string, password: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthState | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const queryClient = useQueryClient()
  const [token, setToken] = useState(() => tokenStore.get())

  // If we have a saved token, ask the server who we are.
  const me = useQuery({
    queryKey: ['me', token],
    queryFn: () => api<User>('/api/auth/me'),
    enabled: !!token,
    retry: false,
  })

  const logout = useCallback(() => {
    tokenStore.clear()
    setToken(null)
    queryClient.clear() // drop cached data from the previous user
  }, [queryClient])

  // Any 401 from the API logs us out.
  useEffect(() => setUnauthorizedHandler(logout), [logout])

  const login = async (email: string, password: string) => {
    const res = await api<LoginResponse>('/api/auth/login', {
      method: 'POST',
      body: { email, password },
    })
    tokenStore.set(res.token)
    queryClient.setQueryData(['me', res.token], res.user)
    setToken(res.token)
  }

  const value: AuthState = {
    user: token ? (me.data ?? null) : null,
    isLoading: !!token && me.isPending,
    login,
    logout,
  }
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used inside <AuthProvider>')
  return ctx
}
