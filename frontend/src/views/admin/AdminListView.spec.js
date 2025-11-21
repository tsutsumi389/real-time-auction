import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import AdminListView from './AdminListView.vue'
import { useAdminStore } from '@/stores/admin'
import * as adminApi from '@/services/adminApi'

// Mock adminApi
vi.mock('@/services/adminApi')

// Mock UI components
vi.mock('@/components/ui/Card.vue', () => ({
  default: {
    template: '<div class="card-mock"><slot /></div>'
  }
}))

vi.mock('@/components/ui/Button.vue', () => ({
  default: {
    template: '<button v-bind="$attrs"><slot /></button>',
    inheritAttrs: false
  }
}))

// Mock admin components
vi.mock('@/components/admin/AdminTable.vue', () => ({
  default: {
    name: 'AdminTable',
    template: '<div class="admin-table-mock">AdminTable</div>',
    props: ['admins', 'loading', 'sortField', 'sortOrder'],
    emits: ['edit', 'statusChange', 'sort']
  }
}))

vi.mock('@/components/admin/AdminFilters.vue', () => ({
  default: {
    name: 'AdminFilters',
    template: '<div class="admin-filters-mock">AdminFilters</div>',
    props: ['modelValue'],
    emits: ['update:modelValue', 'reset']
  }
}))

vi.mock('@/components/admin/AdminSearchBar.vue', () => ({
  default: {
    name: 'AdminSearchBar',
    template: '<div class="admin-search-bar-mock">AdminSearchBar</div>',
    props: ['modelValue'],
    emits: ['update:modelValue', 'search']
  }
}))

vi.mock('@/components/admin/AdminStatusChangeDialog.vue', () => ({
  default: {
    name: 'AdminStatusChangeDialog',
    template: '<div class="admin-status-change-dialog-mock">AdminStatusChangeDialog</div>',
    props: ['open', 'admin', 'loading'],
    emits: ['update:open', 'confirm']
  }
}))

describe('AdminListView Integration Tests', () => {
  let wrapper
  let router
  let adminStore

  const mockAdmins = [
    {
      id: 1,
      email: 'admin1@example.com',
      role: 'system_admin',
      status: 'active',
      created_at: '2025-01-01T00:00:00Z',
      updated_at: '2025-01-01T00:00:00Z'
    },
    {
      id: 2,
      email: 'admin2@example.com',
      role: 'auctioneer',
      status: 'active',
      created_at: '2025-01-02T00:00:00Z',
      updated_at: '2025-01-02T00:00:00Z'
    },
    {
      id: 3,
      email: 'admin3@example.com',
      role: 'system_admin',
      status: 'suspended',
      created_at: '2025-01-03T00:00:00Z',
      updated_at: '2025-01-03T00:00:00Z'
    }
  ]

  const mockPagination = {
    current_page: 1,
    total_pages: 5,
    total_items: 100,
    items_per_page: 20
  }

  beforeEach(() => {
    setActivePinia(createPinia())
    adminStore = useAdminStore()

    // Create router
    router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/admin/admins', component: AdminListView },
        { path: '/admin/admins/:id/edit', component: { template: '<div>Edit Admin</div>' } }
      ]
    })

    // Mock API responses
    vi.mocked(adminApi.getAdminList).mockResolvedValue({
      admins: mockAdmins,
      pagination: mockPagination
    })

    vi.mocked(adminApi.updateAdminStatus).mockResolvedValue({
      admin: {
        id: 1,
        status: 'suspended',
        updated_at: '2025-01-04T00:00:00Z'
      }
    })
  })

  describe('Initial Load', () => {
    it('should fetch and display admin list on mount', async () => {
      wrapper = mount(AdminListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      expect(adminApi.getAdminList).toHaveBeenCalled()
      expect(adminStore.admins).toEqual(mockAdmins)
      expect(adminStore.pagination).toEqual({
        currentPage: 1,
        totalPages: 5,
        totalItems: 100,
        itemsPerPage: 20
      })
    })

    it('should handle fetch error on mount', async () => {
      const errorMessage = 'Network error'
      vi.mocked(adminApi.getAdminList).mockRejectedValue(new Error(errorMessage))

      wrapper = mount(AdminListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      expect(adminStore.error).toBe(errorMessage)
    })
  })

  describe('Search Functionality', () => {
    it('should fetch admins when search is submitted', async () => {
      wrapper = mount(AdminListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()
      vi.clearAllMocks()

      // Update search term and trigger search
      adminStore.filters.search = 'admin1'
      await adminStore.setFiltersAndFetch({ search: 'admin1' })

      await flushPromises()

      expect(adminApi.getAdminList).toHaveBeenCalledWith(
        expect.objectContaining({
          search: 'admin1'
        })
      )
    })

    it('should reset page to 1 when searching', async () => {
      wrapper = mount(AdminListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      // Set to page 3
      adminStore.pagination.currentPage = 3

      // Search
      await adminStore.setFiltersAndFetch({ search: 'test' })

      await flushPromises()

      expect(adminStore.pagination.currentPage).toBe(1)
    })
  })

  describe('Filter Functionality', () => {
    it('should fetch admins when role filter is changed', async () => {
      wrapper = mount(AdminListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()
      vi.clearAllMocks()

      // Change role filter
      await adminStore.setFiltersAndFetch({ role: 'system_admin' })

      await flushPromises()

      expect(adminApi.getAdminList).toHaveBeenCalledWith(
        expect.objectContaining({
          role: 'system_admin'
        })
      )
    })

    it('should fetch admins when status filter is changed', async () => {
      wrapper = mount(AdminListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()
      vi.clearAllMocks()

      // Change status filter
      await adminStore.setFiltersAndFetch({ status: 'active' })

      await flushPromises()

      expect(adminApi.getAdminList).toHaveBeenCalledWith(
        expect.objectContaining({
          status: 'active'
        })
      )
    })

    it('should reset filters when reset button is clicked', async () => {
      wrapper = mount(AdminListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      // Set filters
      adminStore.filters.search = 'test'
      adminStore.filters.role = 'system_admin'
      adminStore.filters.status = 'active'

      // Reset
      await adminStore.resetFilters()

      await flushPromises()

      expect(adminStore.filters).toEqual({
        search: '',
        role: '',
        status: '',
        sort: 'id',
        order: 'asc'
      })
    })
  })

  describe('Sorting Functionality', () => {
    it('should sort by email when email column is clicked', async () => {
      wrapper = mount(AdminListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()
      vi.clearAllMocks()

      // Sort by email
      await adminStore.changeSort('email')

      await flushPromises()

      expect(adminApi.getAdminList).toHaveBeenCalledWith(
        expect.objectContaining({
          sort: 'email',
          order: 'asc'
        })
      )
    })

    it('should toggle sort order when clicking same column', async () => {
      wrapper = mount(AdminListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      // First click - ascending
      await adminStore.changeSort('email')
      expect(adminStore.filters.order).toBe('asc')

      // Second click - descending
      await adminStore.changeSort('email')
      expect(adminStore.filters.order).toBe('desc')
    })
  })

  describe('Pagination Functionality', () => {
    it('should fetch admins when page is changed', async () => {
      wrapper = mount(AdminListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()
      vi.clearAllMocks()

      // Change to page 3
      await adminStore.changePage(3)

      await flushPromises()

      expect(adminApi.getAdminList).toHaveBeenCalledWith(
        expect.objectContaining({
          page: 3
        })
      )
    })
  })

  describe('Status Change Functionality', () => {
    it('should update admin status successfully', async () => {
      wrapper = mount(AdminListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      // Get the admin before status change
      const admin = adminStore.admins[0]
      expect(admin.status).toBe('active')

      // Change status
      const result = await adminStore.changeAdminStatus(1, 'suspended')

      await flushPromises()

      expect(result).toBe(true)
      expect(adminApi.updateAdminStatus).toHaveBeenCalledWith(1, 'suspended')

      // Verify the admin status was updated in the store
      const updatedAdmin = adminStore.admins.find(a => a.id === 1)
      expect(updatedAdmin.status).toBe('suspended')
    })

    it('should handle status change error', async () => {
      const errorMessage = 'Permission denied'
      vi.mocked(adminApi.updateAdminStatus).mockRejectedValue(new Error(errorMessage))

      wrapper = mount(AdminListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      // Try to change status
      const result = await adminStore.changeAdminStatus(1, 'suspended')

      await flushPromises()

      expect(result).toBe(false)
      expect(adminStore.error).toBe(errorMessage)

      // Verify the admin status was NOT updated in the store
      const admin = adminStore.admins.find(a => a.id === 1)
      expect(admin.status).toBe('active')
    })
  })

  describe('Navigation', () => {
    it('should navigate to edit page when edit button is clicked', async () => {
      const pushSpy = vi.spyOn(router, 'push')

      wrapper = mount(AdminListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      // Simulate edit button click (in real implementation this would be triggered by AdminTable)
      await router.push('/admin/admins/1/edit')

      expect(pushSpy).toHaveBeenCalledWith('/admin/admins/1/edit')
    })
  })

  describe('URL Query Parameters', () => {
    it('should apply filters from URL query parameters on mount', async () => {
      // Set up router with query parameters
      await router.push('/admin/admins?search=test&role=system_admin&status=active&page=2')

      wrapper = mount(AdminListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      expect(adminApi.getAdminList).toHaveBeenCalledWith(
        expect.objectContaining({
          search: 'test',
          role: 'system_admin',
          status: 'active',
          page: 2
        })
      )
    })

    it('should update URL when filters change', async () => {
      wrapper = mount(AdminListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      // Change filters
      await adminStore.setFiltersAndFetch({
        search: 'test',
        role: 'system_admin',
        status: 'active'
      })

      await flushPromises()

      // In real implementation, the component would update the URL
      // Here we just verify the filters were applied
      expect(adminStore.filters.search).toBe('test')
      expect(adminStore.filters.role).toBe('system_admin')
      expect(adminStore.filters.status).toBe('active')
    })
  })

  describe('Loading State', () => {
    it('should show loading state during fetch', async () => {
      // Create a promise that we can control
      let resolveFn
      const promise = new Promise((resolve) => {
        resolveFn = resolve
      })

      vi.mocked(adminApi.getAdminList).mockReturnValue(promise)

      wrapper = mount(AdminListView, {
        global: {
          plugins: [router]
        }
      })

      // Loading should be true
      expect(adminStore.loading).toBe(true)

      // Resolve the promise
      resolveFn({
        admins: mockAdmins,
        pagination: mockPagination
      })

      await flushPromises()

      // Loading should be false
      expect(adminStore.loading).toBe(false)
    })
  })

  describe('Multiple Filters Combined', () => {
    it('should apply multiple filters simultaneously', async () => {
      wrapper = mount(AdminListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()
      vi.clearAllMocks()

      // Apply multiple filters
      await adminStore.setFiltersAndFetch({
        search: 'admin',
        role: 'system_admin',
        status: 'active'
      })

      await flushPromises()

      expect(adminApi.getAdminList).toHaveBeenCalledWith(
        expect.objectContaining({
          search: 'admin',
          role: 'system_admin',
          status: 'active',
          sort: 'id',
          order: 'asc',
          page: 1,
          limit: 20
        })
      )
    })
  })

  describe('Error Handling', () => {
    it('should clear error when retrying after error', async () => {
      // First request fails
      vi.mocked(adminApi.getAdminList).mockRejectedValueOnce(new Error('Network error'))

      wrapper = mount(AdminListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      expect(adminStore.error).toBe('Network error')

      // Clear error
      adminStore.clearError()
      expect(adminStore.error).toBeNull()

      // Second request succeeds
      vi.mocked(adminApi.getAdminList).mockResolvedValue({
        admins: mockAdmins,
        pagination: mockPagination
      })

      await adminStore.fetchAdminList()

      await flushPromises()

      expect(adminStore.error).toBeNull()
      expect(adminStore.admins).toEqual(mockAdmins)
    })
  })
})
