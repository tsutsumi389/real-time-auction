import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import AdminRegisterView from './AdminRegisterView.vue'
import * as adminApi from '@/services/adminApi'

// Mock components to avoid dependency issues
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

vi.mock('@/components/ui/Separator.vue', () => ({
  default: {
    template: '<hr class="separator-mock" />'
  }
}))

vi.mock('@/components/ui/RadioGroup.vue', () => ({
  default: {
    template: '<div class="radio-group-mock"><slot /></div>',
    props: ['modelValue', 'name'],
    emits: ['update:modelValue']
  }
}))

vi.mock('@/components/ui/RadioGroupItem.vue', () => ({
  default: {
    template: `
      <div class="radio-item-mock">
        <input type="radio" :id="id" :value="value" @change="$emit('change', value)" />
        <slot />
      </div>
    `,
    props: ['id', 'value'],
    emits: ['change']
  }
}))

describe('AdminRegisterView', () => {
  let wrapper
  let router

  beforeEach(() => {
    // Create a router instance
    router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/admin/register', name: 'admin-register', component: AdminRegisterView },
        { path: '/admin/list', name: 'admin-list', component: { template: '<div>Admin List</div>' } }
      ]
    })

    // Mount component
    wrapper = mount(AdminRegisterView, {
      global: {
        plugins: [router]
      }
    })
  })

  it('renders registration form correctly', () => {
    expect(wrapper.find('h1').text()).toBe('新規管理者登録')
    expect(wrapper.find('#email').exists()).toBe(true)
    expect(wrapper.find('#display_name').exists()).toBe(true)
    expect(wrapper.find('#password').exists()).toBe(true)
    expect(wrapper.find('#password_confirm').exists()).toBe(true)
    expect(wrapper.find('#role-system-admin').exists()).toBe(true)
    expect(wrapper.find('#role-auctioneer').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('displays section headers', () => {
    const headers = wrapper.findAll('h2')
    expect(headers[0].text()).toBe('基本情報')
    expect(headers[1].text()).toBe('認証情報')
    expect(headers[2].text()).toBe('権限設定')
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

  it('does not show error for optional display_name field when empty', async () => {
    const displayNameInput = wrapper.find('#display_name')
    await displayNameInput.trigger('blur')
    // display_name is optional, so no validation error should appear
    const errorElements = wrapper.findAll('.text-red-600')
    const hasDisplayNameError = errorElements.some(el =>
      el.text().includes('表示名') && !el.text().includes('任意')
    )
    expect(hasDisplayNameError).toBe(false)
  })

  it('displays validation error when display_name exceeds max length', async () => {
    const displayNameInput = wrapper.find('#display_name')
    await displayNameInput.setValue('a'.repeat(101))
    await displayNameInput.trigger('blur')
    expect(wrapper.text()).toContain('表示名は100文字以内で入力してください')
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

  it('displays validation error when password confirmation does not match', async () => {
    const passwordInput = wrapper.find('#password')
    const passwordConfirmInput = wrapper.find('#password_confirm')

    await passwordInput.setValue('password123')
    await passwordConfirmInput.setValue('different123')
    await passwordConfirmInput.trigger('blur')

    expect(wrapper.text()).toContain('パスワードが一致しません')
  })

  it('re-validates password confirmation when password changes', async () => {
    const passwordInput = wrapper.find('#password')
    const passwordConfirmInput = wrapper.find('#password_confirm')

    // Set matching passwords
    await passwordInput.setValue('password123')
    await passwordConfirmInput.setValue('password123')
    await passwordConfirmInput.trigger('blur')

    // Change password to make them not match
    await passwordInput.setValue('newpassword123')
    await passwordInput.trigger('blur')

    expect(wrapper.text()).toContain('パスワードが一致しません')
  })

  it('does not submit form when validation fails', async () => {
    const registerAdminSpy = vi.spyOn(adminApi, 'registerAdmin')

    // Submit form without filling in required fields
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()

    expect(registerAdminSpy).not.toHaveBeenCalled()
    expect(wrapper.text()).toContain('入力内容に誤りがあります')
  })

  it('submits form with valid data for system_admin role', async () => {
    const registerAdminSpy = vi.spyOn(adminApi, 'registerAdmin').mockResolvedValue({
      admin: { id: 1, email: 'admin@example.com' }
    })
    const pushSpy = vi.spyOn(router, 'push')

    // Fill in valid data
    await wrapper.find('#email').setValue('admin@example.com')
    await wrapper.find('#display_name').setValue('システム管理者')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#password_confirm').setValue('password123')

    // Select role
    const systemAdminRadio = wrapper.find('#role-system-admin')
    await systemAdminRadio.trigger('change')
    wrapper.vm.formData.role = 'system_admin'

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(registerAdminSpy).toHaveBeenCalledWith({
      email: 'admin@example.com',
      password: 'password123',
      display_name: 'システム管理者',
      role: 'system_admin'
    })
    expect(pushSpy).toHaveBeenCalledWith({ name: 'admin-list' })
  })

  it('submits form with valid data for auctioneer role', async () => {
    const registerAdminSpy = vi.spyOn(adminApi, 'registerAdmin').mockResolvedValue({
      admin: { id: 2, email: 'auctioneer@example.com' }
    })
    const pushSpy = vi.spyOn(router, 'push')

    // Fill in valid data
    await wrapper.find('#email').setValue('auctioneer@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#password_confirm').setValue('password123')

    // Select role
    const auctioneerRadio = wrapper.find('#role-auctioneer')
    await auctioneerRadio.trigger('change')
    wrapper.vm.formData.role = 'auctioneer'

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(registerAdminSpy).toHaveBeenCalledWith({
      email: 'auctioneer@example.com',
      password: 'password123',
      display_name: undefined,
      role: 'auctioneer'
    })
    expect(pushSpy).toHaveBeenCalledWith({ name: 'admin-list' })
  })

  it('displays error message when email already exists', async () => {
    vi.spyOn(adminApi, 'registerAdmin').mockRejectedValue({
      status: 409,
      message: 'Email already exists'
    })

    // Fill in valid data
    await wrapper.find('#email').setValue('existing@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#password_confirm').setValue('password123')
    wrapper.vm.formData.role = 'system_admin'

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(wrapper.text()).toContain('このメールアドレスは既に登録されています')
    expect(wrapper.text()).toContain('メールアドレスが既に使用されています')
  })

  it('displays error message when validation fails on server', async () => {
    vi.spyOn(adminApi, 'registerAdmin').mockRejectedValue({
      status: 400,
      message: 'Validation error'
    })

    // Fill in data
    await wrapper.find('#email').setValue('admin@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#password_confirm').setValue('password123')
    wrapper.vm.formData.role = 'system_admin'

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(wrapper.text()).toContain('Validation error')
  })

  it('displays error message when user lacks permission', async () => {
    vi.spyOn(adminApi, 'registerAdmin').mockRejectedValue({
      status: 403,
      message: 'Forbidden'
    })

    // Fill in data
    await wrapper.find('#email').setValue('admin@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#password_confirm').setValue('password123')
    wrapper.vm.formData.role = 'system_admin'

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(wrapper.text()).toContain('この操作を行う権限がありません')
  })

  it('displays generic error message on server error', async () => {
    vi.spyOn(adminApi, 'registerAdmin').mockRejectedValue({
      status: 500,
      message: 'Internal Server Error'
    })

    // Fill in data
    await wrapper.find('#email').setValue('admin@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#password_confirm').setValue('password123')
    wrapper.vm.formData.role = 'system_admin'

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(wrapper.text()).toContain('サーバーエラーが発生しました')
  })

  it('shows loading state during registration', async () => {
    vi.spyOn(adminApi, 'registerAdmin').mockImplementation(() => {
      return new Promise(resolve => setTimeout(() => resolve({ admin: {} }), 100))
    })

    // Fill in data
    await wrapper.find('#email').setValue('admin@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#password_confirm').setValue('password123')
    wrapper.vm.formData.role = 'system_admin'

    // Submit form
    wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()

    // Check loading state
    expect(wrapper.find('button[type="submit"]').text()).toContain('登録中...')
    expect(wrapper.find('button[type="submit"]').attributes('disabled')).toBeDefined()
  })

  it('disables all inputs during loading', async () => {
    vi.spyOn(adminApi, 'registerAdmin').mockImplementation(() => {
      return new Promise(resolve => setTimeout(() => resolve({ admin: {} }), 100))
    })

    // Fill in data
    await wrapper.find('#email').setValue('admin@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#password_confirm').setValue('password123')
    wrapper.vm.formData.role = 'system_admin'

    // Submit form
    wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()

    // Check inputs are disabled
    expect(wrapper.find('#email').attributes('disabled')).toBeDefined()
    expect(wrapper.find('#display_name').attributes('disabled')).toBeDefined()
    expect(wrapper.find('#password').attributes('disabled')).toBeDefined()
    expect(wrapper.find('#password_confirm').attributes('disabled')).toBeDefined()
  })

  it('disables cancel button during loading', async () => {
    vi.spyOn(adminApi, 'registerAdmin').mockImplementation(() => {
      return new Promise(resolve => setTimeout(() => resolve({ admin: {} }), 100))
    })

    // Fill in data
    await wrapper.find('#email').setValue('admin@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#password_confirm').setValue('password123')
    wrapper.vm.formData.role = 'system_admin'

    // Submit form
    wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()

    // Check cancel button is disabled
    const buttons = wrapper.findAll('button')
    const cancelButton = buttons.find(btn => btn.text().includes('キャンセル'))
    expect(cancelButton.attributes('disabled')).toBeDefined()
  })

  it('navigates back to admin list when cancel button is clicked', async () => {
    const pushSpy = vi.spyOn(router, 'push')

    const buttons = wrapper.findAll('button')
    const cancelButton = buttons.find(btn => btn.text().includes('キャンセル'))

    await cancelButton.trigger('click')

    expect(pushSpy).toHaveBeenCalledWith({ name: 'admin-list' })
  })

  it('navigates back to admin list when back link is clicked', async () => {
    const pushSpy = vi.spyOn(router, 'push')

    const backButton = wrapper.find('button')
    await backButton.trigger('click')

    expect(pushSpy).toHaveBeenCalledWith({ name: 'admin-list' })
  })

  it('clears form error when submitting after previous error', async () => {
    // First submission with error
    vi.spyOn(adminApi, 'registerAdmin').mockRejectedValueOnce({
      status: 500,
      message: 'Server error'
    })

    await wrapper.find('#email').setValue('admin@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#password_confirm').setValue('password123')
    wrapper.vm.formData.role = 'system_admin'

    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(wrapper.text()).toContain('サーバーエラーが発生しました')

    // Second submission success
    vi.spyOn(adminApi, 'registerAdmin').mockResolvedValueOnce({
      admin: { id: 1, email: 'admin@example.com' }
    })

    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(wrapper.text()).not.toContain('サーバーエラーが発生しました')
  })

  it('displays role descriptions correctly', () => {
    expect(wrapper.text()).toContain('システム管理者（system_admin）')
    expect(wrapper.text()).toContain('全権限。ユーザー管理、ポイント管理が可能。')
    expect(wrapper.text()).toContain('主催者（auctioneer）')
    expect(wrapper.text()).toContain('オークション管理のみ。ユーザー管理は不可。')
  })

  it('displays helper texts for inputs', () => {
    expect(wrapper.text()).toContain('任意。未入力の場合はメールアドレスが使用されます')
    expect(wrapper.text()).toContain('8文字以上で入力してください')
    expect(wrapper.text()).toContain('確認のため、もう一度入力してください')
  })

  it('displays required field indicators', () => {
    const requiredIndicators = wrapper.findAll('.text-red-500')
    expect(requiredIndicators.length).toBeGreaterThan(0)
    expect(requiredIndicators.some(el => el.text() === '*')).toBe(true)
  })
})
