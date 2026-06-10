const BASE = import.meta.env.VITE_API_URL || '/api'

async function request(path, options = {}) {
  const res = await fetch(`${BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => null)
    throw new Error(err?.error || `Ошибка ${res.status}`)
  }
  if (res.status === 204) return null
  return res.json()
}

// Categories
export const getCategories = () => request('/categories')
export const createCategory = (data) => request('/categories', { method: 'POST', body: JSON.stringify(data) })
export const updateCategory = (id, data) => request(`/categories/${id}`, { method: 'PUT', body: JSON.stringify(data) })
export const deleteCategory = (id) => request(`/categories/${id}`, { method: 'DELETE' })

// Brands
export const getBrands = () => request('/brands')
export const createBrand = (data) => request('/brands', { method: 'POST', body: JSON.stringify(data) })
export const updateBrand = (id, data) => request(`/brands/${id}`, { method: 'PUT', body: JSON.stringify(data) })
export const deleteBrand = (id) => request(`/brands/${id}`, { method: 'DELETE' })

// Parts
export const getParts = () => request('/parts')
export const getPart = (id) => request(`/parts/${id}`)
export const createPart = (data) => request('/parts', { method: 'POST', body: JSON.stringify(data) })
export const updatePart = (id, data) => request(`/parts/${id}`, { method: 'PUT', body: JSON.stringify(data) })
export const deletePart = (id) => request(`/parts/${id}`, { method: 'DELETE' })

// Stock
export const getStock = () => request('/stock')

// Income
export const getIncome = () => request('/income')
export const createIncome = (data) => request('/income', { method: 'POST', body: JSON.stringify(data) })

// Outcome
export const getOutcome = () => request('/outcome')
export const createOutcome = (data) => request('/outcome', { method: 'POST', body: JSON.stringify(data) })
