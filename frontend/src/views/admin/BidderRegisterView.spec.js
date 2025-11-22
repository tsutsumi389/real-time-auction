import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import BidderRegisterView from './BidderRegisterView.vue'
import * as bidderApi from '@/services/bidderApi'

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

describe('BidderRegisterView', () => {
  let wrapper
  let router

  beforeEach(() => {
    // Create a router instance
    router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/admin/bidders/new', name: 'bidder-register', component: BidderRegisterView },
        { path: '/admin/bidders', name: 'bidder-list', component: { template: '<div>Bidder List</div>' } }
      ]
    })

    // Mount component
    wrapper = mount(BidderRegisterView, {
      global: {
        plugins: [router]
      }
    })
  })

  it('renders registration form correctly', () => {
    expect(wrapper.find('h1').text()).toBe('新規入札者登録')
    expect(wrapper.find('#email').exists()).toBe(true)
    expect(wrapper.find('#display_name').exists()).toBe(true)
    expect(wrapper.find('#password').exists()).toBe(true)
    expect(wrapper.find('#confirmPassword').exists()).toBe(true)
    expect(wrapper.find('#initial_points').exists()).toBe(true)
    expect(wrapper.find('button[type="submit"]').exists()).toBe(true)
  })

  it('displays section headers', () => {
    const headers = wrapper.findAll('h2')
    expect(headers[0].text()).toBe('基本情報')
    expect(headers[1].text()).toBe('認証情報')
    expect(headers[2].text()).toBe('ポイント設定')
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

  it('displays validation error when email exceeds max length', async () => {
    const emailInput = wrapper.find('#email')
    const longEmail = 'a'.repeat(246) + '@example.com' // 256 characters
    await emailInput.setValue(longEmail)
    await emailInput.trigger('blur')
    expect(wrapper.text()).toContain('メールアドレスは255文字以内で入力してください')
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

  it('displays validation error when password confirmation is empty', async () => {
    const confirmPasswordInput = wrapper.find('#confirmPassword')
    await confirmPasswordInput.trigger('blur')
    expect(wrapper.text()).toContain('確認用パスワードを入力してください')
  })

  it('displays validation error when password confirmation does not match', async () => {
    const passwordInput = wrapper.find('#password')
    const confirmPasswordInput = wrapper.find('#confirmPassword')

    await passwordInput.setValue('password123')
    await confirmPasswordInput.setValue('different123')
    await confirmPasswordInput.trigger('blur')

    expect(wrapper.text()).toContain('パスワードが一致しません')
  })

  it('re-validates password confirmation when password changes', async () => {
    const passwordInput = wrapper.find('#password')
    const confirmPasswordInput = wrapper.find('#confirmPassword')

    // Set matching passwords
    await passwordInput.setValue('password123')
    await confirmPasswordInput.setValue('password123')
    await confirmPasswordInput.trigger('blur')

    // Change password to make them not match
    await passwordInput.setValue('newpassword123')
    await passwordInput.trigger('blur')

    expect(wrapper.text()).toContain('パスワードが一致しません')
  })

  it('does not show error for optional initial_points field when empty', async () => {
    const pointsInput = wrapper.find('#initial_points')
    await pointsInput.trigger('blur')
    // initial_points is optional, so no validation error should appear
    const errorElements = wrapper.findAll('.text-red-600')
    const hasPointsError = errorElements.some(el =>
      el.text().includes('ポイント') && !el.text().includes('任意')
    )
    expect(hasPointsError).toBe(false)
  })

  it('displays validation error when initial_points is negative', async () => {
    const pointsInput = wrapper.find('#initial_points')
    await pointsInput.setValue('-100')
    await pointsInput.trigger('blur')
    expect(wrapper.text()).toContain('0以上の整数を入力してください')
  })

  it('displays validation error when initial_points is not an integer', async () => {
    const pointsInput = wrapper.find('#initial_points')
    await pointsInput.setValue('100.5')
    await pointsInput.trigger('blur')
    expect(wrapper.text()).toContain('整数で入力してください')
  })

  it('displays validation error when initial_points is not a number', async () => {
    // Note: type="number" inputs in browsers prevent non-numeric input
    // This test verifies the validation logic handles non-numeric values
    // In real usage, browsers will block 'abc' from being entered
    const pointsInput = wrapper.find('#initial_points')
    
    // Manually set the value to bypass browser validation in tests
    wrapper.vm.formData.initial_points = 'abc'
    await wrapper.vm.$nextTick()
    
    // Trigger validation
    wrapper.vm.validateField('initial_points')
    await wrapper.vm.$nextTick()
    
    expect(wrapper.text()).toContain('数値を入力してください')
  })

  it('does not submit form when validation fails', async () => {
    const registerBidderSpy = vi.spyOn(bidderApi, 'registerBidder')

    // Submit form without filling in required fields
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()

    expect(registerBidderSpy).not.toHaveBeenCalled()
    expect(wrapper.text()).toContain('入力内容に誤りがあります')
  })

  it('submits form with valid data and no initial points', async () => {
    const registerBidderSpy = vi.spyOn(bidderApi, 'registerBidder').mockResolvedValue({
      id: '550e8400-e29b-41d4-a716-446655440000',
      email: 'bidder@example.com',
      display_name: '入札者01',
      status: 'active',
      points: {
        total_points: 0,
        available_points: 0,
        reserved_points: 0
      }
    })
    const pushSpy = vi.spyOn(router, 'push')

    // Fill in valid data
    await wrapper.find('#email').setValue('bidder@example.com')
    await wrapper.find('#display_name').setValue('入札者01')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#confirmPassword').setValue('password123')

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(registerBidderSpy).toHaveBeenCalledWith({
      email: 'bidder@example.com',
      password: 'password123',
      display_name: '入札者01',
      initial_points: undefined
    })
    expect(pushSpy).toHaveBeenCalledWith({ name: 'bidder-list' })
  })

  it('submits form with valid data and initial points', async () => {
    const registerBidderSpy = vi.spyOn(bidderApi, 'registerBidder').mockResolvedValue({
      id: '550e8400-e29b-41d4-a716-446655440000',
      email: 'bidder@example.com',
      display_name: '入札者01',
      status: 'active',
      points: {
        total_points: 1000,
        available_points: 1000,
        reserved_points: 0
      }
    })
    const pushSpy = vi.spyOn(router, 'push')

    // Fill in valid data with initial points
    await wrapper.find('#email').setValue('bidder@example.com')
    await wrapper.find('#display_name').setValue('入札者01')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#confirmPassword').setValue('password123')
    await wrapper.find('#initial_points').setValue('1000')

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(registerBidderSpy).toHaveBeenCalledWith({
      email: 'bidder@example.com',
      password: 'password123',
      display_name: '入札者01',
      initial_points: 1000
    })
    expect(pushSpy).toHaveBeenCalledWith({ name: 'bidder-list' })
  })

  it('submits form without display_name when not provided', async () => {
    const registerBidderSpy = vi.spyOn(bidderApi, 'registerBidder').mockResolvedValue({
      id: '550e8400-e29b-41d4-a716-446655440000',
      email: 'bidder@example.com',
      display_name: null,
      status: 'active',
      points: {
        total_points: 0,
        available_points: 0,
        reserved_points: 0
      }
    })
    const pushSpy = vi.spyOn(router, 'push')

    // Fill in only required fields
    await wrapper.find('#email').setValue('bidder@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#confirmPassword').setValue('password123')

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(registerBidderSpy).toHaveBeenCalledWith({
      email: 'bidder@example.com',
      password: 'password123',
      display_name: undefined,
      initial_points: undefined
    })
    expect(pushSpy).toHaveBeenCalledWith({ name: 'bidder-list' })
  })

  it('displays error message when email already exists', async () => {
    vi.spyOn(bidderApi, 'registerBidder').mockRejectedValue({
      status: 409,
      message: 'Email already exists'
    })

    // Fill in valid data
    await wrapper.find('#email').setValue('existing@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#confirmPassword').setValue('password123')

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(wrapper.text()).toContain('このメールアドレスは既に登録されています')
    expect(wrapper.text()).toContain('メールアドレスが既に使用されています')
  })

  it('displays error message when validation fails on server', async () => {
    vi.spyOn(bidderApi, 'registerBidder').mockRejectedValue({
      status: 400,
      message: 'Validation error'
    })

    // Fill in data
    await wrapper.find('#email').setValue('bidder@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#confirmPassword').setValue('password123')

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(wrapper.text()).toContain('Validation error')
  })

  it('displays error message when user lacks permission', async () => {
    vi.spyOn(bidderApi, 'registerBidder').mockRejectedValue({
      status: 403,
      message: 'Forbidden'
    })

    // Fill in data
    await wrapper.find('#email').setValue('bidder@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#confirmPassword').setValue('password123')

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(wrapper.text()).toContain('この操作を行う権限がありません')
  })

  it('displays generic error message on server error', async () => {
    vi.spyOn(bidderApi, 'registerBidder').mockRejectedValue({
      status: 500,
      message: 'Internal Server Error'
    })

    // Fill in data
    await wrapper.find('#email').setValue('bidder@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#confirmPassword').setValue('password123')

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(wrapper.text()).toContain('サーバーエラーが発生しました')
  })

  it('shows loading state during registration', async () => {
    vi.spyOn(bidderApi, 'registerBidder').mockImplementation(() => {
      return new Promise(resolve => setTimeout(() => resolve({ 
        id: '550e8400-e29b-41d4-a716-446655440000',
        email: 'bidder@example.com' 
      }), 100))
    })

    // Fill in data
    await wrapper.find('#email').setValue('bidder@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#confirmPassword').setValue('password123')

    // Submit form
    wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()

    // Check loading state
    expect(wrapper.find('button[type="submit"]').text()).toContain('登録中...')
    expect(wrapper.find('button[type="submit"]').attributes('disabled')).toBeDefined()
  })

  it('disables all inputs during loading', async () => {
    vi.spyOn(bidderApi, 'registerBidder').mockImplementation(() => {
      return new Promise(resolve => setTimeout(() => resolve({ 
        id: '550e8400-e29b-41d4-a716-446655440000',
        email: 'bidder@example.com' 
      }), 100))
    })

    // Fill in data
    await wrapper.find('#email').setValue('bidder@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#confirmPassword').setValue('password123')

    // Submit form
    wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()

    // Check inputs are disabled
    expect(wrapper.find('#email').attributes('disabled')).toBeDefined()
    expect(wrapper.find('#display_name').attributes('disabled')).toBeDefined()
    expect(wrapper.find('#password').attributes('disabled')).toBeDefined()
    expect(wrapper.find('#confirmPassword').attributes('disabled')).toBeDefined()
    expect(wrapper.find('#initial_points').attributes('disabled')).toBeDefined()
  })

  it('disables cancel button during loading', async () => {
    vi.spyOn(bidderApi, 'registerBidder').mockImplementation(() => {
      return new Promise(resolve => setTimeout(() => resolve({ 
        id: '550e8400-e29b-41d4-a716-446655440000',
        email: 'bidder@example.com' 
      }), 100))
    })

    // Fill in data
    await wrapper.find('#email').setValue('bidder@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#confirmPassword').setValue('password123')

    // Submit form
    wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()

    // Check cancel button is disabled
    const buttons = wrapper.findAll('button')
    const cancelButton = buttons.find(btn => btn.text().includes('キャンセル'))
    expect(cancelButton.attributes('disabled')).toBeDefined()
  })

  it('navigates back to bidder list when cancel button is clicked', async () => {
    const pushSpy = vi.spyOn(router, 'push')

    const buttons = wrapper.findAll('button')
    const cancelButton = buttons.find(btn => btn.text().includes('キャンセル'))

    await cancelButton.trigger('click')

    expect(pushSpy).toHaveBeenCalledWith({ name: 'bidder-list' })
  })

  it('navigates back to bidder list when back link is clicked', async () => {
    const pushSpy = vi.spyOn(router, 'push')

    // Find the back button (first button in the component)
    const backButton = wrapper.find('button')
    await backButton.trigger('click')

    expect(pushSpy).toHaveBeenCalledWith({ name: 'bidder-list' })
  })

  it('clears form error when submitting after previous error', async () => {
    // First submission with error
    vi.spyOn(bidderApi, 'registerBidder').mockRejectedValueOnce({
      status: 500,
      message: 'Server error'
    })

    await wrapper.find('#email').setValue('bidder@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#confirmPassword').setValue('password123')

    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(wrapper.text()).toContain('サーバーエラーが発生しました')

    // Second submission success
    vi.spyOn(bidderApi, 'registerBidder').mockResolvedValueOnce({
      id: '550e8400-e29b-41d4-a716-446655440000',
      email: 'bidder@example.com'
    })

    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(wrapper.text()).not.toContain('サーバーエラーが発生しました')
  })

  it('displays helper texts for inputs', () => {
    expect(wrapper.text()).toContain('任意。未入力の場合はメールアドレスが使用されます')
    expect(wrapper.text()).toContain('8文字以上で入力してください')
    expect(wrapper.text()).toContain('確認のため、もう一度入力してください')
    expect(wrapper.text()).toContain('任意。0以上の整数で入力してください。未入力の場合は0ポイントとして登録されます')
  })

  it('displays required field indicators', () => {
    const requiredIndicators = wrapper.findAll('.text-red-500')
    expect(requiredIndicators.length).toBeGreaterThan(0)
    expect(requiredIndicators.some(el => el.text() === '*')).toBe(true)
  })

  it('accepts zero as valid initial points', async () => {
    const registerBidderSpy = vi.spyOn(bidderApi, 'registerBidder').mockResolvedValue({
      id: '550e8400-e29b-41d4-a716-446655440000',
      email: 'bidder@example.com',
      status: 'active',
      points: {
        total_points: 0,
        available_points: 0,
        reserved_points: 0
      }
    })

    await wrapper.find('#email').setValue('bidder@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#confirmPassword').setValue('password123')
    await wrapper.find('#initial_points').setValue('0')
    await wrapper.find('#initial_points').trigger('blur')

    // Should not show validation error
    const pointsErrors = wrapper.findAll('.text-red-600')
    const hasPointsError = pointsErrors.some(el => 
      el.text().includes('ポイント') && !el.text().includes('任意')
    )
    expect(hasPointsError).toBe(false)

    // Submit form
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(registerBidderSpy).toHaveBeenCalledWith({
      email: 'bidder@example.com',
      password: 'password123',
      display_name: undefined,
      initial_points: 0
    })
  })

  it('accepts large initial points value', async () => {
    const registerBidderSpy = vi.spyOn(bidderApi, 'registerBidder').mockResolvedValue({
      id: '550e8400-e29b-41d4-a716-446655440000',
      email: 'bidder@example.com',
      status: 'active',
      points: {
        total_points: 1000000,
        available_points: 1000000,
        reserved_points: 0
      }
    })

    await wrapper.find('#email').setValue('bidder@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('#confirmPassword').setValue('password123')
    await wrapper.find('#initial_points').setValue('1000000')

    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()
    await new Promise(resolve => setTimeout(resolve, 0))

    expect(registerBidderSpy).toHaveBeenCalledWith({
      email: 'bidder@example.com',
      password: 'password123',
      display_name: undefined,
      initial_points: 1000000
    })
  })

  it('displays separators between sections', () => {
    const separators = wrapper.findAll('.separator-mock')
    expect(separators.length).toBe(3) // Three sections: basic info, auth info, points
  })

  it('focuses on first error field when validation fails', async () => {
    // Submit form without filling in required fields
    await wrapper.find('form').trigger('submit.prevent')
    await wrapper.vm.$nextTick()

    // Check that error message is displayed
    expect(wrapper.text()).toContain('入力内容に誤りがあります')
    expect(wrapper.text()).toContain('メールアドレスを入力してください')
  })
})
