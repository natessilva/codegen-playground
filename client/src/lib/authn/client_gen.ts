// Code generated - DO NOT EDIT.

import { getCaller } from "../fetch";
 
export class AuthNService {
  constructor(private readonly url: string, private readonly token: string) {}

  private readonly caller = getCaller(this.url, "authn", "AuthNService", this.token);
 	
	exchangeEmailToken = async (input:ExchangeEmailTokenInput):Promise<AuthOutput> =>
	  this.caller("ExchangeEmailToken", input);
	
	login = async (input:AuthInput):Promise<AuthOutput> =>
	  this.caller("Login", input);
	
	resetPassword = async (input:ResetPasswordInput):Promise<ResetPasswordOutput> =>
	  this.caller("ResetPassword", input);
	
	signup = async (input:AuthInput):Promise<AuthOutput> =>
	  this.caller("Signup", input);
	
}



export type ExchangeEmailTokenInput = {
	token:string
}

export type AuthOutput = {
	token:string
	ok:boolean
}

export type AuthInput = {
	email:string
	password:string
}

export type ResetPasswordInput = {
	email:string
}

export type ResetPasswordOutput = {
	ok:boolean
}

