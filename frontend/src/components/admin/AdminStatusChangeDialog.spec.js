import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import AdminStatusChangeDialog from './AdminStatusChangeDialog.vue'

describe('AdminStatusChangeDialog', () => {
  const activeAdmin = {
    id: 1,
    email: 'admin@example.com',
    role: 'system_admin',
    status: 'active',
  }

  const inactiveAdmin = {
    id: 2,
    email: 'inactive@example.com',
    role: 'auctioneer',
    status: 'inactive',
  }

  describe('Rendering', () => {
    it('should not render when modelValue is false', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: false,
          admin: activeAdmin,
        },
      })

      expect(wrapper.find('[role="dialog"]').exists()).toBe(false)
    })

    it('should render when modelValue is true', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: activeAdmin,
        },
      })

      expect(wrapper.find('[role="dialog"]').exists()).toBe(true)
    })

    it('should display admin email', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: activeAdmin,
        },
      })

      expect(wrapper.text()).toContain('admin@example.com')
    })
  })

  describe('Active Admin (Deactivation)', () => {
    it('should show correct title for deactivation', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: activeAdmin,
        },
      })

      expect(wrapper.find('#modal-title').text()).toBe('アカウント停止の確認')
    })

    it('should show warning message for deactivation', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: activeAdmin,
        },
      })

      expect(wrapper.text()).toContain('停止後、このアカウントではログインできなくなります。')
    })

    it('should show red warning icon for deactivation', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: activeAdmin,
        },
      })

      const iconContainer = wrapper.find('.flex-shrink-0')
      expect(iconContainer.classes()).toContain('bg-red-100')
      expect(iconContainer.find('.text-red-600').exists()).toBe(true)
    })

    it('should show correct button text for deactivation', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: activeAdmin,
        },
      })

      const confirmButton = wrapper.findAll('button')[0]
      expect(confirmButton.text()).toBe('停止する')
    })

    it('should have red button for deactivation', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: activeAdmin,
        },
      })

      const confirmButton = wrapper.findAll('button')[0]
      expect(confirmButton.classes()).toContain('bg-red-600')
    })
  })

  describe('Inactive Admin (Activation)', () => {
    it('should show correct title for activation', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: inactiveAdmin,
        },
      })

      expect(wrapper.find('#modal-title').text()).toBe('アカウント復活の確認')
    })

    it('should not show warning message for activation', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: inactiveAdmin,
        },
      })

      expect(wrapper.text()).not.toContain('停止後、このアカウントではログインできなくなります。')
    })

    it('should show green success icon for activation', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: inactiveAdmin,
        },
      })

      const iconContainer = wrapper.find('.flex-shrink-0')
      expect(iconContainer.classes()).toContain('bg-green-100')
      expect(iconContainer.find('.text-green-600').exists()).toBe(true)
    })

    it('should show correct button text for activation', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: inactiveAdmin,
        },
      })

      const confirmButton = wrapper.findAll('button')[0]
      expect(confirmButton.text()).toBe('復活する')
    })

    it('should have green button for activation', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: inactiveAdmin,
        },
      })

      const confirmButton = wrapper.findAll('button')[0]
      expect(confirmButton.classes()).toContain('bg-green-600')
    })
  })

  describe('User Interactions', () => {
    it('should emit confirm event when clicking confirm button', async () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: activeAdmin,
        },
      })

      const confirmButton = wrapper.findAll('button')[0]
      await confirmButton.trigger('click')

      expect(wrapper.emitted('confirm')).toBeTruthy()
    })

    it('should emit update:modelValue with false when clicking cancel button', async () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: activeAdmin,
        },
      })

      const cancelButton = wrapper.findAll('button')[1]
      await cancelButton.trigger('click')

      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')[0]).toEqual([false])
    })

    it('should emit update:modelValue when clicking backdrop', async () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: activeAdmin,
        },
      })

      const backdrop = wrapper.find('.bg-gray-500')
      await backdrop.trigger('click')

      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')[0]).toEqual([false])
    })
  })

  describe('Loading State', () => {
    it('should show loading spinner when loading', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: activeAdmin,
          loading: true,
        },
      })

      const confirmButton = wrapper.findAll('button')[0]
      expect(confirmButton.find('.animate-spin').exists()).toBe(true)
    })

    it('should disable confirm button when loading', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: activeAdmin,
          loading: true,
        },
      })

      const confirmButton = wrapper.findAll('button')[0]
      expect(confirmButton.attributes('disabled')).toBeDefined()
      expect(confirmButton.classes()).toContain('disabled:opacity-50')
    })

    it('should disable cancel button when loading', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: activeAdmin,
          loading: true,
        },
      })

      const cancelButton = wrapper.findAll('button')[1]
      expect(cancelButton.attributes('disabled')).toBeDefined()
    })

    it('should not close dialog when loading', async () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: activeAdmin,
          loading: true,
        },
      })

      const cancelButton = wrapper.findAll('button')[1]
      await cancelButton.trigger('click')

      // loadingがtrueの場合、update:modelValueは発行されない
      expect(wrapper.emitted('update:modelValue')).toBeFalsy()
    })

    it('should not close dialog via backdrop when loading', async () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: activeAdmin,
          loading: true,
        },
      })

      const backdrop = wrapper.find('.bg-gray-500')
      await backdrop.trigger('click')

      expect(wrapper.emitted('update:modelValue')).toBeFalsy()
    })
  })

  describe('Accessibility', () => {
    it('should have proper ARIA attributes', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: activeAdmin,
        },
      })

      const dialog = wrapper.find('[role="dialog"]')
      expect(dialog.attributes('aria-modal')).toBe('true')
      expect(dialog.attributes('aria-labelledby')).toBe('modal-title')
    })

    it('should have modal title with id', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: activeAdmin,
        },
      })

      const title = wrapper.find('#modal-title')
      expect(title.exists()).toBe(true)
    })
  })

  describe('Edge Cases', () => {
    it('should handle null admin gracefully', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: null,
        },
      })

      // エラーなくレンダリングされる
      expect(wrapper.find('[role="dialog"]').exists()).toBe(true)
    })

    it('should handle undefined admin gracefully', () => {
      const wrapper = mount(AdminStatusChangeDialog, {
        props: {
          modelValue: true,
          admin: undefined,
        },
      })

      // エラーなくレンダリングされる
      expect(wrapper.find('[role="dialog"]').exists()).toBe(true)
    })
  })
})
