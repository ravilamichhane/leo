// Code generated by tygo. DO NOT EDIT.

import { Fetcher } from "./fetcher";
import { SafeParse } from "./safeparse";

//////////
// source: CreateUser.go

export interface CreateUserRequest {
  Name: string;
  Age: number /* int */;
}
export interface CreateUserRespone {}
export const createUser = async (
  data: CreateUserRequest,
  id: string,
  options?: any,
) => {
  return await SafeParse(
    Fetcher<CreateUserRespone>({
      url: `/User/${id}`,
      body: data,
      method: "POST",
      headers: {
        ...options?.headers,
      },
      ...options,
    }),
  );
};

//////////
// source: DeleteUser.go

export interface DeleeteUserRespone {}
export const deleteUser = async (id: string, options?: any) => {
  return await SafeParse(
    Fetcher<DeleeteUserRespone>({
      url: `/User/${id}`,
      method: "DELETE",
      headers: {
        ...options?.headers,
      },
      ...options,
    }),
  );
};

//////////
// source: GetAllUser.go

export interface GetAllUsersRespone {}
export const getAllUsers = async (options?: any) => {
  return await SafeParse(
    Fetcher<Page<GetAllUsersRespone>>({
      url: `/User`,
      method: "Get",
      headers: {
        ...options?.headers,
      },
      ...options,
    }),
  );
};

//////////
// source: GetUserById.go

export interface GetUserByIdRespone {}
export const getUserById = async (id: string, options?: any) => {
  return await SafeParse(
    Fetcher<GetUserByIdRespone>({
      url: `/User/${id}`,
      method: "Get",
      headers: {
        ...options?.headers,
      },
      ...options,
    }),
  );
};

//////////
// source: UpdateUser.go

export interface UpdateUserRequest {}
export interface UpdateUserRespone {}
export const updateUser = async (
  data: UpdateUserRequest,
  id: string,
  options?: any,
) => {
  return await SafeParse(
    Fetcher<UpdateUserRespone>({
      url: `/User/${id}`,
      body: data,
      method: "PUT",
      headers: {
        ...options?.headers,
      },
      ...options,
    }),
  );
};