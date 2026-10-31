import { useAuth } from '../auth/AuthContext'

// Placeholder — expiry and low-stock cards arrive in Phase 3.
export function DashboardPage() {
  const { user } = useAuth()
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-800">Welcome, {user?.name}</h1>
      <p className="mt-2 text-gray-600">You are signed in as <strong>{user?.role}</strong>.</p>
    </div>
  )
}
