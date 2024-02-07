export type TypeSigninBody = {
  email: string;
  password: string;
}

export type TypeSignupBody = {
  first_name: string;
  last_name: string;
  email: string;
  password: string;
  captcha: string;
}

export interface APIResponse {
  message?: string;
  token?: string;
  error?: string;
  error_description?: string;
}

export interface TypeGetClientResponse {
  id: string;
  description: string;
  message: string;
  name: string;
  redirect_uri: string;
  response_type: string;
  scope: Array<string>;
}

export interface TypeAuthGrant {
  response_type: string | null;
  client_id: string | null;
  redirect_uri: string | null;
  state: string | null;
}

export interface TypeClient {
  id: string;
  name: string;
  redirect_uri: string;
  response_type: string;
  scope: Array<string>;
}
