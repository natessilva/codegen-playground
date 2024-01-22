// Code generated - DO NOT EDIT.

import { getCaller } from "../fetch";
 
export class SpaceService {
  constructor(private readonly url: string, private readonly token: string) {}

  private readonly caller = getCaller(this.url, "app", "SpaceService", this.token);
 	
	get = async (input:Empty):Promise<Space> =>
	  this.caller("Get", input);
	
	getUsers = async (input:Empty):Promise<GetUsersResponse> =>
	  this.caller("GetUsers", input);
	
	list = async (input:Empty):Promise<ListSpacesResponse> =>
	  this.caller("List", input);
	
	update = async (input:Space):Promise<Empty> =>
	  this.caller("Update", input);
	
}
 
export class TicketService {
  constructor(private readonly url: string, private readonly token: string) {}

  private readonly caller = getCaller(this.url, "app", "TicketService", this.token);
 	
	assign = async (input:AssignInput):Promise<OK> =>
	  this.caller("Assign", input);
	
	assignSelf = async (input:ID):Promise<OK> =>
	  this.caller("AssignSelf", input);
	
	create = async (input:Ticket):Promise<ID> =>
	  this.caller("Create", input);
	
	get = async (input:ID):Promise<GetTicketResponse> =>
	  this.caller("Get", input);
	
	update = async (input:UpdateTicketInput):Promise<OK> =>
	  this.caller("Update", input);
	
}
 
export class UserService {
  constructor(private readonly url: string, private readonly token: string) {}

  private readonly caller = getCaller(this.url, "app", "UserService", this.token);
 	
	get = async (input:Empty):Promise<User> =>
	  this.caller("Get", input);
	
	setPassword = async (input:SetPasswordInput):Promise<OK> =>
	  this.caller("SetPassword", input);
	
	update = async (input:UpdateUserInput):Promise<Empty> =>
	  this.caller("Update", input);
	
}



export type UpdateUserInput = {
	name:string
}

export type ListSpacesResponse = {
	spaces:Space[]
}

export type AssignInput = {
	ticketID:number
	userID:number
}

export type UpdateTicketInput = {
	id:number
	subject:string
	body:string
}

export type Empty = {
}

export type OK = {
	ok:boolean
}

export type Ticket = {
	subject:string
	body:string
}

export type GetUsersResponse = {
	users:User[]
}

export type ID = {
	id:number
}

export type GetTicketResponse = {
	ticket:Ticket
	ok:boolean
}

export type SetPasswordInput = {
	oldPassword:string
	newPassword:string
}

export type Space = {
	name:string
}

export type User = {
	id:number
	name:string
}

