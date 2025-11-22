import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import BidderStatusChangeDialog from './BidderStatusChangeDialog.vue'

describe('BidderStatusChangeDialog', () => {
  const activeBidder = {
    id: 'abc12345-def6-7890-abcd-ef1234567890',
    email: 'bidder1@example.com',
    display_name: '田中太郎',
    status: 'active',
  }

  const suspendedBidder = {
    id: 'def67890-abc1-2345-6789-0abcdef12345',
    email: 'bidder2@example.com',
    display_name: '佐藤花子',
    status: 'suspended',
  }

  describe('Rendering', () => {
    it('should not render when open is false', () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: false,
          bidder: activeBidder,
        },
      })

      expect(wrapper.find('[role="dialog"]').exists()).toBe(false)
    })

    it('should render when open is true', () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: activeBidder,
        },
      })

      expect(wrapper.find('[role="dialog"]').exists()).toBe(true)
    })

    it('should display bidder email', () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: activeBidder,
        },
      })

      expect(wrapper.text()).toContain('bidder1@example.com')
    })

    it('should display bidder display name', () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: activeBidder,
        },
      })

      expect(wrapper.text()).toContain('田中太郎')
    })
  })

  describe('Suspend Action (active -> suspended)', () => {
    it('should show correct title for suspension', () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: activeBidder,
        },
      })

      expect(wrapper.find('#modal-title').text()).toBe('アカウント停止の確認')
    })

    it('should show warning message for suspension', () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: activeBidder,
        },
      })

      expect(wrapper.text()).toContain('停止後、このアカウントではログインできなくなります。')
    })

    it('should show red warning icon for suspension', () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: activeBidder,
        },
      })

      const iconContainer = wrapper.find('.flex-shrink-0')
      expect(iconContainer.classes()).toContain('bg-red-100')
      expect(iconContainer.find('.text-red-600').exists()).toBe(true)
    })

    it('should show correct button text for suspension', () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: activeBidder,
        },
      })

      const confirmButton = wrapper.findAll('button')[0]
      expect(confirmButton.text()).toBe('停止する')
    })

    it('should have red button for suspension', () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: activeBidder,
        },
      })

      const confirmButton = wrapper.findAll('button')[0]
      expect(confirmButton.classes()).toContain('bg-red-600')
    })
  })

  describe('Activate Action (suspended -> active)', () => {
    it('should show correct title for activation', () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: suspendedBidder,
        },
      })

      expect(wrapper.find('#modal-title').text()).toBe('アカウント復活の確認')
    })

    it('should not show warning message for activation', () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: suspendedBidder,
        },
      })

      expect(wrapper.text()).not.toContain('停止後、このアカウントではログインできなくなります。')
    })

    it('should show green success icon for activation', () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: suspendedBidder,
        },
      })

      const iconContainer = wrapper.find('.flex-shrink-0')
      expect(iconContainer.classes()).toContain('bg-green-100')
      expect(iconContainer.find('.text-green-600').exists()).toBe(true)
    })

    it('should show correct button text for activation', () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: suspendedBidder,
        },
      })

      const confirmButton = wrapper.findAll('button')[0]
      expect(confirmButton.text()).toBe('復活する')
    })

    it('should have green button for activation', () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: suspendedBidder,
        },
      })

      const confirmButton = wrapper.findAll('button')[0]
      expect(confirmButton.classes()).toContain('bg-green-600')
    })
  })

  describe('User Interactions', () => {
    it('should emit confirm event when clicking confirm button', async () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: activeBidder,
        },
      })

      const confirmButton = wrapper.findAll('button')[0]
      await confirmButton.trigger('click')

      expect(wrapper.emitted('confirm')).toBeTruthy()
    })

    it('should emit update:modelValue with false when clicking cancel button', async () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: activeBidder,
        },
      })

      const cancelButton = wrapper.findAll('button')[1]
      await cancelButton.trigger('click')

      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')[0]).toEqual([false])
    })

    it('should emit update:modelValue when clicking backdrop', async () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: activeBidder,
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
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: activeBidder,
          loading: true,
        },
      })

      const confirmButton = wrapper.findAll('button')[0]
      expect(confirmButton.find('.animate-spin').exists()).toBe(true)
    })

    it('should disable confirm button when loading', () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: activeBidder,
          loading: true,
        },
      })

      const confirmButton = wrapper.findAll('button')[0]
      expect(confirmButton.attributes('disabled')).toBeDefined()
      expect(confirmButton.classes()).toContain('disabled:opacity-50')
    })

    it('should disable cancel button when loading', () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: activeBidder,
          loading: true,
        },
      })

      const cancelButton = wrapper.findAll('button')[1]
      expect(cancelButton.attributes('disabled')).toBeDefined()
    })

    it('should not close dialog when loading', async () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: activeBidder,
          loading: true,
        },
      })

      const cancelButton = wrapper.findAll('button')[1]
      await cancelButton.trigger('click')

      // loadingがtrueの場合、update:modelValueは発行されない
      expect(wrapper.emitted('update:modelValue')).toBeFalsy()
    })

    it('should not close dialog via backdrop when loading', async () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: activeBidder,
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
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: activeBidder,
        },
      })

      const dialog = wrapper.find('[role="dialog"]')
      expect(dialog.attributes('aria-modal')).toBe('true')
      expect(dialog.attributes('aria-labelledby')).toBe('modal-title')
    })

    it('should have modal title with id', () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: activeBidder,
        },
      })

      const title = wrapper.find('#modal-title')
      expect(title.exists()).toBe(true)
    })
  })

  describe('Edge Cases', () => {
    it('should handle null bidder gracefully', () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: null,
        },
      })

      // エラーなくレンダリングされる
      expect(wrapper.find('[role="dialog"]').exists()).toBe(true)
    })

    it('should handle undefined bidder gracefully', () => {
      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: undefined,
        },
      })

      // エラーなくレンダリングされる
      expect(wrapper.find('[role="dialog"]').exists()).toBe(true)
    })

    it('should handle bidder with null display_name', () => {
      const bidderWithoutName = {
        ...activeBidder,
        display_name: null,
      }

      const wrapper = mount(BidderStatusChangeDialog, {
        props: {
          modelValue: true,
          bidder: bidderWithoutName,
        },
      })

      // display_nameがnullでもエラーなくレンダリングされる
      expect(wrapper.find('[role="dialog"]').exists()).toBe(true)
    })
  })
})
