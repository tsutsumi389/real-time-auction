import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import GrantPointsDialog from './GrantPointsDialog.vue'

// Mock Dialog component
const mockDialogComponent = {
  template: '<div v-if="modelValue" class=[role="dialog"]"><slot /></div>',
  props: ['modelValue']
}

describe('GrantPointsDialog', () => {
  const mockBidder = {
    id: 'abc12345-def6-7890-abcd-ef1234567890',
    email: 'bidder1@example.com',
    display_name: '田中太郎',
    points: 10000,
  }

  describe('Rendering', () => {
    it('should not render when closed', () => {
      const wrapper = mount(GrantPointsDialog, {
        props: {
          modelValue: false,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
          },
        },
      })

      expect(wrapper.find('[role="dialog"]').exists()).toBe(false)
    })

    it('should render when open', () => {
      const wrapper = mount(GrantPointsDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
          },
        },
      })

      expect(wrapper.find('[role="dialog"]').exists()).toBe(true)
    })

    it('should display bidder information', () => {
      const wrapper = mount(GrantPointsDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
          },
        },
      })

      expect(wrapper.text()).toContain('bidder1@example.com')
      expect(wrapper.text()).toContain('田中太郎')
      expect(wrapper.text()).toContain('10,000')
    })

    it('should display input field for points', () => {
      const wrapper = mount(GrantPointsDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
          },
        },
      })

      const input = wrapper.find('input[type="number"]')
      expect(input.exists()).toBe(true)
    })

    it('should display preview of total points after grant', async () => {
      const wrapper = mount(GrantPointsDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
          },
        },
      })

      const input = wrapper.find('input[type="number"]')
      await input.setValue(1000)

      expect(wrapper.text()).toContain('11,000')
    })
  })

  describe('Validation', () => {
    it('should show error for negative points', async () => {
      const wrapper = mount(GrantPointsDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
          },
        },
      })

      const input = wrapper.find('input[type="number"]')
      await input.setValue(-100)

      expect(wrapper.text()).toContain('ポイントは1以上の整数を入力してください')
    })

    it('should show error for zero points', async () => {
      const wrapper = mount(GrantPointsDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
          },
        },
      })

      const input = wrapper.find('input[type="number"]')
      await input.setValue(0)

      expect(wrapper.text()).toContain('ポイントは1以上の整数を入力してください')
    })

    it('should disable submit button when points are invalid', async () => {
      const wrapper = mount(GrantPointsDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
          },
        },
      })

      const input = wrapper.find('input[type="number"]')
      await input.setValue(-100)

      const submitButton = wrapper.findAll('button').find((btn) => btn.text().includes('付与'))
      expect(submitButton.attributes('disabled')).toBeDefined()
    })

    it('should enable submit button when points are valid', async () => {
      const wrapper = mount(GrantPointsDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
          },
        },
      })

      const input = wrapper.find('input[type="number"]')
      await input.setValue(1000)

      const submitButton = wrapper.findAll('button').find((btn) => btn.text().includes('付与'))
      expect(submitButton.attributes('disabled')).toBeUndefined()
    })
  })

  describe('Actions', () => {
    it('should emit update:modelValue event when cancel button is clicked', async () => {
      const wrapper = mount(GrantPointsDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
          },
        },
      })

      const cancelButton = wrapper.findAll('button').find((btn) => btn.text().includes('キャンセル'))
      await cancelButton.trigger('click')

      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')[0]).toEqual([false])
    })

    it('should emit confirm event with points when submit button is clicked', async () => {
      const wrapper = mount(GrantPointsDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
          },
        },
      })

      const input = wrapper.find('input[type="number"]')
      await input.setValue(1000)

      const submitButton = wrapper.findAll('button').find((btn) => btn.text().includes('付与'))
      await submitButton.trigger('click')

      expect(wrapper.emitted('confirm')).toBeTruthy()
      expect(wrapper.emitted('confirm')[0]).toEqual([1000])
    })

    it('should disable buttons when loading', async () => {
      const wrapper = mount(GrantPointsDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
          loading: true,
        },
        global: {
          stubs: {
            
          },
        },
      })

      const buttons = wrapper.findAll('button')
      buttons.forEach((button) => {
        expect(button.attributes('disabled')).toBeDefined()
      })
    })

    it('should clear input when dialog is closed', async () => {
      const wrapper = mount(GrantPointsDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
          },
        },
      })

      const input = wrapper.find('input[type="number"]')
      await input.setValue(1000)

      await wrapper.setProps({ modelValue: false })
      await wrapper.setProps({ modelValue: true })

      expect(input.element.value).toBe('')
    })
  })

  describe('Points Formatting', () => {
    it('should format current points with comma separator', () => {
      const wrapper = mount(GrantPointsDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
          },
        },
      })

      expect(wrapper.text()).toContain('10,000')
    })

    it('should format preview points with comma separator', async () => {
      const wrapper = mount(GrantPointsDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
          },
        },
      })

      const input = wrapper.find('input[type="number"]')
      await input.setValue(5000)

      expect(wrapper.text()).toContain('15,000')
    })
  })
})
