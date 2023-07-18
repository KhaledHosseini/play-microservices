'use client'

import { useAuth } from "@/components/providers/auth_provider";
import SignUpForm from "@/components/client/signup_form";
import { redirect } from 'next/navigation'

export default function SignUp() {
  
  const [user,_setUser] = useAuth()
  if (user) {
      redirect("/")
  }
  
  return (
    <main>
      <SignUpForm/>
    </main>
  );
}