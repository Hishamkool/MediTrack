import { Link } from 'react-router-dom'

export function NotFoundPage() {
  return (
    <div className="py-16 text-center">
      <h1 className="text-2xl font-bold text-gray-800">Page not found</h1>
      <Link to="/" className="mt-4 inline-block text-teal-700 underline">
        Go to dashboard
      </Link>
    </div>
  )
}
