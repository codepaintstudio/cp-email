import request from '@/utils/request'

interface SendEmailRequest {
  email: string;
  password: string;
  subject: string;
  receiverItemsArray: Array<(string | { [k: string]: any })[]>;
  content: string;
}
export const sendEmailService = (data: SendEmailRequest) => request.post('/email/sendemail', data)

