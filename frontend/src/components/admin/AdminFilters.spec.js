import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import AdminFilters from './AdminFilters.vue'

describe('AdminFilters', () => {
  const defaultModelValue = {
    search: '',
    role: '',
    status: '',
  }

  describe('Rendering', () => {
    it('should render all filter inputs', () => {
      const wrapper = mount(AdminFilters, {
        props: {
          modelValue: defaultModelValue,
        },
      })

      expect(wrapper.find('#search').exists()).toBe(true)
      expect(wrapper.find('#role').exists()).toBe(true)
      expect(wrapper.find('#status').exists()).toBe(true)
    })

    it('should render action buttons', () => {
      const wrapper = mount(AdminFilters, {
        props: {
          modelValue: defaultModelValue,
        },
      })

      const buttons = wrapper.findAll('button')
      expect(buttons).toHaveLength(2)
      expect(buttons[0].text()).toBe('検索')
      expect(buttons[1].text()).toBe('リセット')
    })

    it('should display current filter values', () => {
      const wrapper = mount(AdminFilters, {
        props: {
          modelValue: {
            search: 'admin@example.com',
            role: 'system_admin',
            status: 'active',
          },
        },
      })

      expect(wrapper.find('#search').element.value).toBe('admin@example.com')
      expect(wrapper.find('#role').element.value).toBe('system_admin')
      expect(wrapper.find('#status').element.value).toBe('active')
    })

    it('should disable buttons when loading', () => {
      const wrapper = mount(AdminFilters, {
        props: {
          modelValue: defaultModelValue,
          loading: true,
        },
      })

      const buttons = wrapper.findAll('button')
      buttons.forEach((button) => {
        expect(button.attributes('disabled')).toBeDefined()
      })
    })
  })

  describe('Search Input', () => {
    it('should emit update:modelValue on search input', async () => {
      const wrapper = mount(AdminFilters, {
        props: {
          modelValue: defaultModelValue,
        },
      })

      const searchInput = wrapper.find('#search')
      await searchInput.setValue('test@example.com')

      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')[0][0]).toEqual({
        search: 'test@example.com',
        role: '',
        status: '',
      })
    })

    it('should emit search event on Enter key', async () => {
      const wrapper = mount(AdminFilters, {
        props: {
          modelValue: defaultModelValue,
        },
      })

      const searchInput = wrapper.find('#search')
      await searchInput.trigger('keyup.enter')

      expect(wrapper.emitted('search')).toBeTruthy()
    })

    it('should not emit filter-change on search input', async () => {
      const wrapper = mount(AdminFilters, {
        props: {
          modelValue: defaultModelValue,
        },
      })

      const searchInput = wrapper.find('#search')
      await searchInput.setValue('test@example.com')

      expect(wrapper.emitted('filter-change')).toBeFalsy()
    })
  })

  describe('Role Filter', () => {
    it('should emit update:modelValue on role change', async () => {
      const wrapper = mount(AdminFilters, {
        props: {
          modelValue: defaultModelValue,
        },
      })

      const roleSelect = wrapper.find('#role')
      await roleSelect.setValue('system_admin')

      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')[0][0]).toEqual({
        search: '',
        role: 'system_admin',
        status: '',
      })
    })

    it('should emit filter-change on role change', async () => {
      const wrapper = mount(AdminFilters, {
        props: {
          modelValue: defaultModelValue,
        },
      })

      const roleSelect = wrapper.find('#role')
      await roleSelect.setValue('system_admin')

      expect(wrapper.emitted('filter-change')).toBeTruthy()
    })

    it('should have all role options', () => {
      const wrapper = mount(AdminFilters, {
        props: {
          modelValue: defaultModelValue,
        },
      })

      const roleSelect = wrapper.find('#role')
      const options = roleSelect.findAll('option')

      expect(options).toHaveLength(3)
      expect(options[0].text()).toBe('すべて')
      expect(options[0].element.value).toBe('')
      expect(options[1].text()).toBe('システム管理者')
      expect(options[1].element.value).toBe('system_admin')
      expect(options[2].text()).toBe('オークショニア')
      expect(options[2].element.value).toBe('auctioneer')
    })
  })

  describe('Status Filter', () => {
    it('should emit update:modelValue on status change', async () => {
      const wrapper = mount(AdminFilters, {
        props: {
          modelValue: defaultModelValue,
        },
      })

      const statusSelect = wrapper.find('#status')
      await statusSelect.setValue('active')

      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')[0][0]).toEqual({
        search: '',
        role: '',
        status: 'active',
      })
    })

    it('should emit filter-change on status change', async () => {
      const wrapper = mount(AdminFilters, {
        props: {
          modelValue: defaultModelValue,
        },
      })

      const statusSelect = wrapper.find('#status')
      await statusSelect.setValue('active')

      expect(wrapper.emitted('filter-change')).toBeTruthy()
    })

    it('should have all status options', () => {
      const wrapper = mount(AdminFilters, {
        props: {
          modelValue: defaultModelValue,
        },
      })

      const statusSelect = wrapper.find('#status')
      const options = statusSelect.findAll('option')

      expect(options).toHaveLength(3)
      expect(options[0].text()).toBe('すべて')
      expect(options[0].element.value).toBe('')
      expect(options[1].text()).toBe('有効')
      expect(options[1].element.value).toBe('active')
      expect(options[2].text()).toBe('停止中')
      expect(options[2].element.value).toBe('inactive')
    })
  })

  describe('Action Buttons', () => {
    it('should emit search event when clicking search button', async () => {
      const wrapper = mount(AdminFilters, {
        props: {
          modelValue: defaultModelValue,
        },
      })

      const searchButton = wrapper.findAll('button')[0]
      await searchButton.trigger('click')

      expect(wrapper.emitted('search')).toBeTruthy()
    })

    it('should emit reset event when clicking reset button', async () => {
      const wrapper = mount(AdminFilters, {
        props: {
          modelValue: defaultModelValue,
        },
      })

      const resetButton = wrapper.findAll('button')[1]
      await resetButton.trigger('click')

      expect(wrapper.emitted('reset')).toBeTruthy()
    })
  })

  describe('Accessibility', () => {
    it('should have proper labels for inputs', () => {
      const wrapper = mount(AdminFilters, {
        props: {
          modelValue: defaultModelValue,
        },
      })

      const searchLabel = wrapper.find('label[for="search"]')
      const roleLabel = wrapper.find('label[for="role"]')
      const statusLabel = wrapper.find('label[for="status"]')

      expect(searchLabel.text()).toBe('メールアドレス検索')
      expect(roleLabel.text()).toBe('ロール')
      expect(statusLabel.text()).toBe('状態')
    })

    it('should have placeholder for search input', () => {
      const wrapper = mount(AdminFilters, {
        props: {
          modelValue: defaultModelValue,
        },
      })

      const searchInput = wrapper.find('#search')
      expect(searchInput.attributes('placeholder')).toBe('例: admin@example.com')
    })
  })

  describe('Complex Scenarios', () => {
    it('should handle multiple filter changes', async () => {
      const wrapper = mount(AdminFilters, {
        props: {
          modelValue: defaultModelValue,
        },
      })

      // 検索入力
      const searchInput = wrapper.find('#search')
      await searchInput.setValue('test@example.com')

      // ロール選択
      const roleSelect = wrapper.find('#role')
      await roleSelect.setValue('system_admin')

      // 状態選択
      const statusSelect = wrapper.find('#status')
      await statusSelect.setValue('active')

      // すべてのupdate:modelValueイベントが発行されている
      expect(wrapper.emitted('update:modelValue')).toHaveLength(3)

      // filter-changeはロールと状態の変更時のみ
      expect(wrapper.emitted('filter-change')).toHaveLength(2)
    })

    it('should maintain other filter values when changing one', async () => {
      const wrapper = mount(AdminFilters, {
        props: {
          modelValue: {
            search: 'test@example.com',
            role: 'system_admin',
            status: '',
          },
        },
      })

      const statusSelect = wrapper.find('#status')
      await statusSelect.setValue('active')

      // 他のフィルタ値が保持されている
      expect(wrapper.emitted('update:modelValue')[0][0]).toEqual({
        search: 'test@example.com',
        role: 'system_admin',
        status: 'active',
      })
    })
  })
})
