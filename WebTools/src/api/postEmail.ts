import request from '@/utils/request'

interface SendEmailRequest {
  email: string;
  password: string;
  subject: string;
  receiverItemsArray: string[][];
  content: string;
}
export const sendEmailService = (data: SendEmailRequest) => request.post('/email/sendemail', data)

