import api from './index'

// Auth
export const login = (data) => api.post('/login', data)
export const loginByQr = (data) => api.post('/loginbyqr', data)
export const refreshToken = (data) => api.post('/refresh',{}, { withCredentials: true })

// User
export const getUsers = (params) => api.get('/view.user', { params })
export const createUser = (data) => api.post('/add.user', data)
export const updateUser = (id, data) => api.put(`/edit.user/${id}`, data)
export const toggleUserStatus = (id) => api.put(`/toggle.status.user/${id}`)
export const deleteuser = (id) => api.delete(`/delete.user/${id}`)
export const changePassword = (data) => api.put('/change.password', data)
export const countuser = () => api.get('/count.user')
export const getrole = () => api.get('/view.role')
export const getuserdata = () => api.get(`/view.user.data`)
export const getuserapprove = () => api.get(`/view.user.approve`)
export const verifyuser = (id) => api.put(`verify.user/${id}`)
export const logoutUser = () => api.post('/logout')

// Backup
export const triggerBackup = () => api.post('add.backup')
export const listBackups = () => api.get('view.backup')
export const downloadBackup = (filename) => api.get('view.download.backup', { 
  params: { file: filename },
  responseType: 'blob' 
})
export const deleteBackup = (filename) => api.delete('delete.backup',{
  params: {file: filename},
})

// RoleHasPermission
export const getrolehaspermission = (id) => api.get(`/view.role.has.permission/${id}`)
export const addrolehaspermission = (data) => api.post('/add.role.has.permission',data)
export const deleterolehaspermission = (data) => api.delete('/delete.role.has.permission',{data})
export const editrole = (id,data) => api.put(`/edit.role/${id}`,data)

// company

export const getcompany = (params) => api.get(`/view.company`,{params})
export const addcompany = (data) => api.post('/add.company',data)
export const updatecompany = (id,data) => api.put(`/update.company/${id}`,data)


// branch
export const addbranch = (data) => api.post('/add.Branch',data)
export const updatebranch = (id,data) => api.put(`/update.Branch/${id}`,data)


// Customer

export const getcustomer = (params) => api.get(`/view.Customer`,{params})
export const addcustomer = (data) => api.post('/add.Customer',data)
export const updatecustomer = (id,data) => api.put(`/update.Customer/${id}`,data)