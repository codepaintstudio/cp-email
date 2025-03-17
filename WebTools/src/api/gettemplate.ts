import request from '../utils/request'
export const getTemplateDownload = () =>
  request.get('/template/xlsx', {
    responseType: 'blob'
  })
