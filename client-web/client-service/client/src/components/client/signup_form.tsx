"use client";
import React from "react";
import { useForm } from "react-hook-form";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from 'next/navigation';
import { CreateUserRequest } from "@/types";

export default function SignUpForm() {
  const router = useRouter();

  const {
    handleSubmit,
    register,
    formState: { errors },
  } = useForm<CreateUserRequest>();
  
  const {
    isLoading,
    isError,
    mutate: createUserMutation,
  } = useMutation((params: CreateUserRequest) =>
    fetch("/api/user/create", {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(params),
    })
  , {
    onSuccess: async (response) => {
      const data = await response.json()
      console.log("response is:",data );
      if (response.ok) {
        router.push('/login');
      }
    },
  });

  const handleCreateUser = (createUserRequest: CreateUserRequest) => {
    createUserRequest.role = Number(createUserRequest.role)
    createUserMutation(createUserRequest);
  };

  return (
    <div className="container mx-auto py-8">
    <h1 className="text-2xl font-bold mb-6 text-center">Sign Up</h1>
    <form className="w-full max-w-sm mx-auto bg-white p-8 rounded-md shadow-md"
    onSubmit={handleSubmit(handleCreateUser)}>
      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="name">Name</label>
        <input className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:border-indigo-500"
          type="text" 
          id="name"
          placeholder="Name"
          {...register('name', { required: true })}/>
      </div>
      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="email">Email</label>
        <input className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:border-indigo-500"
          type="email" 
          id="email"
          placeholder="example@example.com"
          {...register('email', { required: true })}/>
      </div>
      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="password">Password</label>
        <input className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:border-indigo-500"
          type="password" 
          id="password"
          placeholder="********"
          {...register('password', { required: true })}/>
      </div>
      <div className="mb-4">
        <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="role">Role: (0: admin)(1: user)</label>
        <input className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:border-indigo-500"
          type="number" 
          id="role"
          placeholder="0"
          {...register('role', { required: true })}/>
      </div>
      <button
        className="w-full bg-indigo-500 text-white text-sm font-bold py-2 px-4 rounded-md hover:bg-indigo-600 transition duration-300"
        type="submit">Sign up</button>
    </form>
  </div>
  );
}