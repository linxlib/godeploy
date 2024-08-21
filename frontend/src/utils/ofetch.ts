import { ofetch } from 'ofetch'
import { showError } from './notification';

export const apiFetch  = ofetch.create({ 
    baseURL: import.meta.env.DEPLOY_BASE_URL,
    credentials: 'include',
    timeout: 10000,
    onResponseError: async ({request,options,response}) =>{
        console.log(
            "[fetch response error]",
            request,
            response.status,
            response.body
          );
          showError(response.statusText)
    }
})

