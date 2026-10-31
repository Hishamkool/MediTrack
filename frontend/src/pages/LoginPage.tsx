import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Navigate, useLocation, useNavigate } from 'react-router-dom'
import { useAuth } from '../auth/AuthContext'
import { FormField, buttonClass, inputClass } from '../components/FormField'

const schema = z.object({
  email: z.string().trim().email('Enter a valid email'),
  password: z.string().min(1, 'Enter your password'),
})
type LoginForm = z.infer<typeof schema>

export function LoginPage() {
  const { user, login } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  // After login, go back to the page the user originally wanted.
  const from = (location.state as { from?: Location } | null)?.from?.pathname ?? '/'

  const {
    register,
    handleSubmit,
    setError,
    formState: { errors, isSubmitting },
  } = useForm<LoginForm>({ resolver: zodResolver(schema) })

  if (user) return <Navigate to={from} replace />

  const onSubmit = async (data: LoginForm) => {
    try {
      await login(data.email, data.password)
      navigate(from, { replace: true })
    } catch (err) {
      setError('root', { message: (err as Error).message })
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4">
      <form onSubmit={handleSubmit(onSubmit)} className="w-full max-w-sm space-y-4 rounded-lg bg-white p-6 shadow" noValidate>
        <div className="text-center">
          <h1 className="text-2xl font-bold text-teal-800">MediTrack</h1>
          <p className="text-sm text-gray-500">Sign in to continue</p>
        </div>

        {errors.root && <div className="rounded bg-red-50 p-3 text-sm text-red-700">{errors.root.message}</div>}

        <FormField label="Email" error={errors.email?.message}>
          <input type="email" autoComplete="username" className={inputClass} {...register('email')} />
        </FormField>
        <FormField label="Password" error={errors.password?.message}>
          <input type="password" autoComplete="current-password" className={inputClass} {...register('password')} />
        </FormField>

        <button type="submit" disabled={isSubmitting} className={`${buttonClass} w-full`}>
          {isSubmitting ? 'Signing in…' : 'Sign in'}
        </button>
      </form>
    </div>
  )
}
