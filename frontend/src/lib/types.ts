// Shapes of the JSON the backend sends. Keep in sync with internal/models.

export type Role = 'admin' | 'staff'

export interface User {
  id: number
  name: string
  email: string
  role: Role
  is_active: boolean
  created_at: string
}

export interface LoginResponse {
  token: string
  expires_at: string
  user: User
}
