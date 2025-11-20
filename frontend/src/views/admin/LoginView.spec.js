import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import LoginView from './LoginView.vue'
import { useAuthStore } from '@/stores/auth'

// Mock components to avoid dependency issues
vi.mock('@/components/ui/Card.vue', () => ({
  default: {
    template: '<div class="card-mock"><slot /></div>'
  }
}))

vi.mock('@/components/ui/Input.vue', () => ({
  default: {
    template: '<input v-bind="$attrs" @input="$emit(\'update:modelValue\', $event.target.value)" />',
    inheritAttrs: false
  }
}))

vi.mock('@/components/ui/Label.vue', () => ({
  default: {
    template: '<label><slot /></label>'
  }
}))

vi.mock('@/components/ui/Button.vue', () => ({
  default: {
    template: '<button v-bind="$attrs"><slot /></button>',
    inheritAttrs: false
  }
}))

vi.mock('@/components/ui/Alert.vue', () => ({
  default: {
    template: '<div class="alert-mock"><slot /></div>'
  }
}))

describe('LoginView', () => {
  let wrapper
  let router
  let authStore

  beforeEach(() => {
    // Create a new pinia instance for each test
    setActivePinia(createPinia())

    // Create a router instance
    router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/admin/login', component: LoginView },
        { path: '/admin/dashboard', component: { template: '<div>Dashboard</div>' } }
      ]
    })

    // Get auth store
    authStore = useAuthStore()

    // Mount component
    wrapper = mount(LoginView, {
      global: {
        plugins: [router]
      }
    })
  })

  it('renders login form correctly', () => {
    expect(wrapper.find('h1').text()).toBe('管理者ログイン')
    expect(wrapper.find('#email').exists()).toBe(true)
    expect(wrapper.find('#password').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('displays validation error when email is empty', async () => {
    const emailInput = wrapper.find('#email')

    await emailInput.trigger('blur')

    expect(wrapper.text()).toContain('メールアドレスを入力してください')
  })

  it('displays validation error when email format is invalid', async () => {
    const emailInput = wrapper.find('#email')

    await emailInput.setValue('invalid-email')
    await emailInput.trigger('blur')

    expect(wrapper.text()).toContain('メールアドレスの形式が正しくありません')
  })

  it('displays validation error when password is empty', async () => {
    const passwordInput = wrapper.find('#password')

    await passwordInput.trigger('blur')

    expect(wrapper.text()).toContain('パスワードを入力してください')
  })

  it('displays validation error when password is too short', async () => {
    const passwordInput = wrapper.find('#password')

    await passwordInput.setValue('short')
    await passwordInput.trigger('blur')

    expect(wrapper.text()).toContain('パスワードは8文字以上で入力してください')
  })

  it('does not submit form when validation fails', async () => {
    const loginSpy = vi.spyOn(authStore, 'login')

    // Submit form without filling in fields
    await wrapper.find('form').trigger('submit.prevent')

    expect(loginSpy).not.toHaveBeenCalled()
  })

  it('submits form with valid credentials', async () => {
    const loginSpy = vi.spyOn(authStore, 'login').mockResolvedValue(true)
    const pushSpy = vi.spyOn(router, 'push')

    // Fill in valid credentials
    await wrapper.find('#email').setValue('admin@example.com')
    await wrapper.find('#password').setValue('password123')

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')

    // Wait for async operations
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(loginSpy).toHaveBeenCalledWith('admin@example.com', 'password123')
    expect(pushSpy).toHaveBeenCalledWith('/admin/dashboard')
  })

  it('displays error message when login fails', async () => {
    const errorMessage = 'メールアドレスまたはパスワードが正しくありません'
    vi.spyOn(authStore, 'login').mockImplementation(async () => {
      authStore.error = errorMessage
      return false
    })

    // Fill in credentials
    await wrapper.find('#email').setValue('admin@example.com')
    await wrapper.find('#password').setValue('wrongpassword')

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')

    // Wait for async operations
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(wrapper.text()).toContain(errorMessage)
  })

  it('shows loading state during login', async () => {
    vi.spyOn(authStore, 'login').mockImplementation(() => {
      return new Promise(resolve => setTimeout(() => resolve(true), 100))
    })

    // Fill in credentials
    await wrapper.find('#email').setValue('admin@example.com')
    await wrapper.find('#password').setValue('password123')

    // Submit form
    wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()

    // Check loading state
    expect(wrapper.find('button[type="submit"]').text()).toBe('ログイン中...')
    expect(wrapper.find('button[type="submit"]').attributes('disabled')).toBeDefined()
  })

  it('disables inputs during loading', async () => {
    vi.spyOn(authStore, 'login').mockImplementation(() => {
      return new Promise(resolve => setTimeout(() => resolve(true), 100))
    })

    // Fill in credentials
    await wrapper.find('#email').setValue('admin@example.com')
    await wrapper.find('#password').setValue('password123')

    // Submit form
    wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()

    // Check inputs are disabled
    expect(wrapper.find('#email').attributes('disabled')).toBeDefined()
    expect(wrapper.find('#password').attributes('disabled')).toBeDefined()
  })

  it('redirects to custom path from query parameter', async () => {
    const loginSpy = vi.spyOn(authStore, 'login').mockResolvedValue(true)
    const pushSpy = vi.spyOn(router, 'push')

    // Navigate to login with redirect query
    await router.push('/admin/login?redirect=/admin/users')

    // Remount component to get new route
    wrapper = mount(LoginView, {
      global: {
        plugins: [router]
      }
    })

    // Fill in credentials and submit
    await wrapper.find('#email').setValue('admin@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('form').trigger('submit.prevent')

    // Wait for async operations
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(pushSpy).toHaveBeenCalledWith('/admin/users')
  })

  it('handles server connection error', async () => {
    vi.spyOn(authStore, 'login').mockRejectedValue(new Error('Network error'))

    // Fill in credentials
    await wrapper.find('#email').setValue('admin@example.com')
    await wrapper.find('#password').setValue('password123')

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')

    // Wait for async operations
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(wrapper.text()).toContain('サーバーに接続できません')
  })
})
