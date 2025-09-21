import request from '../utils/request'
export const uploadPhotoService = (data: FormData) => request.post('/uploads/file', data)
