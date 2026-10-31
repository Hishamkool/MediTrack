import { useState } from 'react'
import { NavLink, Outlet } from 'react-router-dom'
import { useAuth } from '../auth/AuthContext'
import type { Role } from '../lib/types'

// One list drives both the desktop and mobile menus. `roles` limits who sees
// a link; leave it out to show it to everyone. Later phases add more items.
const NAV_ITEMS: { to: string; label: string; roles?: Role[] }[] = [
  { to: '/', label: 'Dashboard' },
  { to: '/users', label: 'Users', roles: ['admin'] },
]

export function Layout() {
  const { user, logout } = useAuth()
  const [menuOpen, setMenuOpen] = useState(false)
  const items = NAV_ITEMS.filter((i) => !i.roles || (user && i.roles.includes(user.role)))

  const linkClass = ({ isActive }: { isActive: boolean }) =>
    `block rounded px-3 py-2 text-sm font-medium ${isActive ? 'bg-teal-700 text-white' : 'text-teal-50 hover:bg-teal-600'}`

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-teal-800">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-3">
          <span className="text-lg font-bold text-white">MediTrack</span>

          {/* Desktop nav */}
          <nav className="hidden items-center gap-1 md:flex">
            {items.map((i) => (
              <NavLink key={i.to} to={i.to} end className={linkClass}>
                {i.label}
              </NavLink>
            ))}
          </nav>
          <div className="hidden items-center gap-3 md:flex">
            <span className="text-sm text-teal-100">
              {user?.name} <span className="rounded bg-teal-600 px-1.5 py-0.5 text-xs uppercase">{user?.role}</span>
            </span>
            <button onClick={logout} className="rounded border border-teal-400 px-3 py-1 text-sm text-white hover:bg-teal-700">
              Log out
            </button>
          </div>

          {/* Mobile menu button */}
          <button
            className="rounded p-2 text-white md:hidden"
            onClick={() => setMenuOpen((o) => !o)}
            aria-label="Toggle menu"
            aria-expanded={menuOpen}
          >
            <svg className="h-6 w-6" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
              <path strokeLinecap="round" d={menuOpen ? 'M6 6l12 12M18 6L6 18' : 'M4 6h16M4 12h16M4 18h16'} />
            </svg>
          </button>
        </div>

        {/* Mobile nav */}
        {menuOpen && (
          <nav className="space-y-1 border-t border-teal-700 px-4 py-3 md:hidden">
            {items.map((i) => (
              <NavLink key={i.to} to={i.to} end className={linkClass} onClick={() => setMenuOpen(false)}>
                {i.label}
              </NavLink>
            ))}
            <div className="flex items-center justify-between pt-2 text-sm text-teal-100">
              <span>
                {user?.name} ({user?.role})
              </span>
              <button onClick={logout} className="rounded border border-teal-400 px-3 py-1 text-white">
                Log out
              </button>
            </div>
          </nav>
        )}
      </header>

      <main className="mx-auto max-w-6xl px-4 py-6">
        <Outlet />
      </main>
    </div>
  )
}
