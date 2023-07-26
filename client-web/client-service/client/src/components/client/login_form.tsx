"use client";

import React from "react";
import { useForm } from "react-hook-form";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from 'next/navigation';
import { LoginUserRequest } from "@/types";
import { toast } from "react-hot-toast";

export default function LogInForm() {
  const router = useRouter()
  const {
    handleSubmit,
    register,
    formState: { errors },
  } = useForm<LoginUserRequest>();
  
  
  const {
    isLoading,
    isError,
    mutate: loginUserMutation,
  } = useMutation((params: LoginUserRequest) =>
    fetch("/api/user/login", {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(params),
    })
  , {
    onSuccess: async (response) => {
      console.log("response is:",response);
      if (response.ok) {
        const data = await response.json()
        console.log("login response is:",data );
        toast.success("Login success.")
        router.push('/');
      }else {
        toast.error("login failed:"+ response.statusText)
      }
    },
    onError: async (error)=> {
      toast.error("login failed:" + error)
    }
  });

  const handleLoginUser = (loginUserRequest: LoginUserRequest) => {
    loginUserMutation(loginUserRequest);
  };

  return (
    <div className="container mx-auto py-8">
    <h1 className="text-2xl font-bold mb-6 text-center">Login</h1>
    <form className="w-full max-w-sm mx-auto bg-white p-8 rounded-md shadow-md"
    onSubmit={handleSubmit(handleLoginUser)}>
      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="email">Email</label>
        <input className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:border-indigo-500"
          type="email" 
          id="Email"
          placeholder="example@example.com"
          {...register('Email', { required: true })}/>
      </div>
      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="password">Password</label>
        <input className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:border-indigo-500"
          type="password" 
          id="Password"
          placeholder="********"
          {...register('Password', { required: true })}/>
      </div>
      <button
        className="w-full bg-indigo-500 text-white text-sm font-bold py-2 px-4 rounded-md hover:bg-indigo-600 transition duration-300"
        disabled={isLoading}
        type="submit">Log in</button>
    </form>
  </div>
  );
}