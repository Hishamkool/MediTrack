// A tiny wrapper around fetch() that:
//  - prefixes the backend URL
//  - attaches the JWT token
//  - turns error responses into an ApiError you can show to the user

const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'
const TOKEN_KEY = 'meditrack_token'

export const tokenStore = {
  get: () => localStorage.getItem(TOKEN_KEY),
  set: (token: string) => localStorage.setItem(TOKEN_KEY, token),
  clear: () => localStorage.removeItem(TOKEN_KEY),
}

export class ApiError extends Error {
  status: number
  fields?: Record<string, string> // per-field validation messages (HTTP 422)

  constructor(status: number, message: string, fields?: Record<string, string>) {
    super(message)
    this.status = status
    this.fields = fields
  }
}

// Called when the server says our token is no longer valid (set by AuthProvider).
let onUnauthorized: () => void = () => {}
export function setUnauthorizedHandler(fn: () => void) {
  onUnauthorized = fn
}

export async function api<T>(path: string, options: { method?: string; body?: unknown } = {}): Promise<T> {
  const headers: Record<string, string> = {}
  const token = tokenStore.get()
  if (token) headers.Authorization = `Bearer ${token}`
  if (options.body !== undefined) headers['Content-Type'] = 'application/json'

  let res: Response
  try {
    res = await fetch(`${API_URL}${path}`, {
      method: options.method ?? 'GET',
      headers,
      body: options.body !== undefined ? JSON.stringify(options.body) : undefined,
    })
  } catch {
    throw new ApiError(0, 'Cannot reach the server. Check your connection.')
  }

  const data = res.status === 204 ? null : await res.json().catch(() => null)

  if (!res.ok) {
    // A 401 on any call except login means the token expired or the user was deactivated.
    if (res.status === 401 && token && path !== '/api/auth/login') onUnauthorized()
    throw new ApiError(res.status, data?.error ?? `Request failed (${res.status})`, data?.fields)
  }
  return data as T
}
