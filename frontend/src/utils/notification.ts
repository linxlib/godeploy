
import { notification, message } from 'ant-design-vue'

export const msgSuccess = (msg:string) => {
    message.success(msg)
}
export const msgError = (msg:string) => {
    message.error(msg)
}
export const msgInfo = (msg:string) => {
  message.info(msg)
}

export const showError = (message: string) => {
    notification.open({
      message: '错误',
      description: message,
      duration: 2,
      type: 'error'
    })
  }
  export const showInfo = (message: string) => {
    notification.open({
      message: '信息',
      description: message,
      duration: 2,
      type: 'info'
    })
  }
  export const showSuccess = (message: string) => {
    notification.open({
      message: '信息',
      description: message,
      duration: 2,
      type: 'success'
    })
  }