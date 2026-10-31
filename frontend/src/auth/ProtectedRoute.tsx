import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { useAuth } from './AuthContext'
import type { Role } from '../lib/types'
import { Spinner } from '../components/Spinner'

// Wrap routes that need a login. Pass `roles` to also restrict by role.
// NOTE: this only hides pages; the backend enforces the real permissions.
export function ProtectedRoute({ roles }: { roles?: Role[] }) {
  const { user, isLoading } = useAuth()
  const location = useLocation()

  if (isLoading) return <Spinner fullScreen />
  if (!user) return <Navigate to="/login" replace state={{ from: location }} />
  if (roles && !roles.includes(user.role)) return <Navigate to="/" replace />
  return <Outlet />
}
