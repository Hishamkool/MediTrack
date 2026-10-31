export function Spinner({ fullScreen = false }: { fullScreen?: boolean }) {
  const spinner = (
    <div className="h-8 w-8 animate-spin rounded-full border-4 border-teal-600 border-t-transparent" role="status" aria-label="Loading" />
  )
  if (!fullScreen) return <div className="flex justify-center p-6">{spinner}</div>
  return <div className="flex min-h-screen items-center justify-center">{spinner}</div>
}
