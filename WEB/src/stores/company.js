import { defineStore } from 'pinia'
import { viewcompanycolor } from '../api/services'

export const useCompanyStore = defineStore('company', {
  state: () => ({
    color: '#1b3351',
  }),
  actions: {
    async fetchColor() {
      try {
       //  const res = await viewcompanycolor()
        this.color =  '#1b3351'
      } catch (e) {
        // keep default on failure
      }
    },
    setColor(color) {
      this.color = color || '#1b3351'
    }
  }
})