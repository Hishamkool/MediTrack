import type { ReactNode } from 'react'

// Label + input + error message, so every form looks the same.
export function FormField({ label, error, children }: { label: string; error?: string; children: ReactNode }) {
  return (
    <label className="block">
      <span className="mb-1 block text-sm font-medium text-gray-700">{label}</span>
      {children}
      {error && <span className="mt-1 block text-sm text-red-600">{error}</span>}
    </label>
  )
}

export const inputClass =
  'w-full rounded border border-gray-300 px-3 py-2 text-base focus:border-teal-600 focus:outline-none focus:ring-1 focus:ring-teal-600'

export const buttonClass =
  'rounded bg-teal-700 px-4 py-2 font-medium text-white hover:bg-teal-800 disabled:cursor-not-allowed disabled:opacity-60'
