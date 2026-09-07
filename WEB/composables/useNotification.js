import { ElNotification } from 'element-plus'

export function useNotification() {
  const success = (message, title = 'Success') => {
    ElNotification({
      title,
      message,
      type: 'success',
      duration: 3000,
      position: 'top-right',
    })
  }

  const error = (message, title = 'Error') => {
    ElNotification({
      title,
      message,
      type: 'error',
      duration: 4000,
      position: 'top-right',
    })
  }

  const warning = (message, title = 'Warning') => {
    ElNotification({
      title,
      message,
      type: 'warning',
      duration: 3000,
      position: 'top-right',
    })
  }

  const info = (message, title = 'Info') => {
    ElNotification({
      title,
      message,
      type: 'info',
      duration: 3000,
      position: 'top-right',
    })
  }

  return {
    success,
    error,
    warning,
    info,
  }
}