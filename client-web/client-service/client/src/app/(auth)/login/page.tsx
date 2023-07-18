'use client'
import { useAuth } from "@/components/providers/auth_provider";
import LogInForm from "@/components/client/login_form";
import { redirect } from 'next/navigation'

export default function LogIn() {
  
  const [user,_setUser] = useAuth()
  if (user) {
      redirect("/")
  }

  return (
    <main>
      <LogInForm/>
    </main>
  );
}