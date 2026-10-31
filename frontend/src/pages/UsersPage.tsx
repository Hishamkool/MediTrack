// Admin-only screen: list users, create staff, activate/deactivate.
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { api, ApiError } from '../lib/api'
import type { User } from '../lib/types'
import { useAuth } from '../auth/AuthContext'
import { Spinner } from '../components/Spinner'
import { FormField, buttonClass, inputClass } from '../components/FormField'

const createSchema = z.object({
  name: z.string().trim().min(1, 'Name is required'),
  email: z.string().trim().email('Enter a valid email'),
  password: z.string().min(8, 'At least 8 characters'),
  role: z.enum(['staff', 'admin']),
})
type CreateUserForm = z.infer<typeof createSchema>

export function UsersPage() {
  const queryClient = useQueryClient()
  const { user: me } = useAuth()

  const users = useQuery({ queryKey: ['users'], queryFn: () => api<User[]>('/api/users') })

  const toggleActive = useMutation({
    mutationFn: (u: User) => api<User>(`/api/users/${u.id}`, { method: 'PATCH', body: { is_active: !u.is_active } }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['users'] }),
    onError: (err) => alert(err.message),
  })

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-gray-800">Users</h1>

      <CreateUserCard />

      <section className="rounded-lg bg-white shadow">
        <h2 className="border-b px-4 py-3 font-semibold text-gray-700">All users</h2>
        {users.isPending && <Spinner />}
        {users.isError && <p className="p-4 text-red-600">{users.error.message}</p>}
        {users.data && (
          <ul className="divide-y">
            {users.data.map((u) => (
              <li key={u.id} className="flex flex-col gap-2 px-4 py-3 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <p className="font-medium text-gray-800">
                    {u.name}{' '}
                    <span className="ml-1 rounded bg-gray-100 px-1.5 py-0.5 text-xs uppercase text-gray-600">{u.role}</span>
                    {!u.is_active && <span className="ml-1 rounded bg-red-100 px-1.5 py-0.5 text-xs text-red-700">inactive</span>}
                  </p>
                  <p className="text-sm text-gray-500">{u.email}</p>
                </div>
                {u.id !== me?.id && (
                  <button
                    onClick={() => {
                      const action = u.is_active ? 'Deactivate' : 'Activate'
                      if (confirm(`${action} ${u.name}?`)) toggleActive.mutate(u)
                    }}
                    disabled={toggleActive.isPending}
                    className={`self-start rounded border px-3 py-1 text-sm sm:self-auto ${
                      u.is_active ? 'border-red-300 text-red-700 hover:bg-red-50' : 'border-teal-300 text-teal-700 hover:bg-teal-50'
                    }`}
                  >
                    {u.is_active ? 'Deactivate' : 'Activate'}
                  </button>
                )}
              </li>
            ))}
          </ul>
        )}
      </section>
    </div>
  )
}

function CreateUserCard() {
  const queryClient = useQueryClient()
  const {
    register,
    handleSubmit,
    reset,
    setError,
    formState: { errors, isSubmitting, isSubmitSuccessful },
  } = useForm<CreateUserForm>({
    resolver: zodResolver(createSchema),
    defaultValues: { role: 'staff' },
  })

  const onSubmit = async (data: CreateUserForm) => {
    try {
      await api<User>('/api/users', { method: 'POST', body: data })
      reset()
      queryClient.invalidateQueries({ queryKey: ['users'] })
    } catch (err) {
      // Show server-side validation messages next to the right field.
      if (err instanceof ApiError && err.fields) {
        for (const [field, msg] of Object.entries(err.fields)) {
          setError(field as keyof CreateUserForm, { message: msg })
        }
      } else {
        setError('root', { message: (err as Error).message })
      }
      throw err // keeps isSubmitSuccessful false
    }
  }

  return (
    <section className="rounded-lg bg-white p-4 shadow">
      <h2 className="mb-3 font-semibold text-gray-700">Add user</h2>
      <form onSubmit={handleSubmit(onSubmit)} noValidate className="grid gap-4 sm:grid-cols-2">
        <FormField label="Name" error={errors.name?.message}>
          <input className={inputClass} {...register('name')} />
        </FormField>
        <FormField label="Email" error={errors.email?.message}>
          <input type="email" className={inputClass} {...register('email')} />
        </FormField>
        <FormField label="Password" error={errors.password?.message}>
          <input type="password" autoComplete="new-password" className={inputClass} {...register('password')} />
        </FormField>
        <FormField label="Role" error={errors.role?.message}>
          <select className={inputClass} {...register('role')}>
            <option value="staff">Staff</option>
            <option value="admin">Admin</option>
          </select>
        </FormField>

        <div className="flex items-center gap-3 sm:col-span-2">
          <button type="submit" disabled={isSubmitting} className={buttonClass}>
            {isSubmitting ? 'Saving…' : 'Create user'}
          </button>
          {isSubmitSuccessful && <span className="text-sm text-teal-700">User created.</span>}
          {errors.root && <span className="text-sm text-red-600">{errors.root.message}</span>}
        </div>
      </form>
    </section>
  )
}
