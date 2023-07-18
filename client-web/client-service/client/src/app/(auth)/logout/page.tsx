'use client'

import React, { useState } from "react";
import { useEffect } from "react";
import { toast } from "react-hot-toast";
import { redirect } from 'next/navigation'

export default function LogOut() {
    const [response,setResponse] = useState<Response | null>(null)
    useEffect(() => {
        // declare the data fetching function
        const fetchData = async () => {
          const response = await fetch("/api/user/logout", {
            method: 'POST'
          });
          console.log("LogOut.useEffect.fetchData: logout result: ",response.body)
          if (response.status == 200){
            toast.success('Logout successfull');
          }else {
            toast.error('Logout failed');
          }
          setResponse(response)
        }
        fetchData()
        .catch(console.error);
      }, [])

    if (response != null) {
        redirect(`/`)
    }

  return (
    <main>
     <h2 className="text-3xl font-bold text-gray-800 dark:text-white md:text-4xl">Logging out</h2>
    </main>
  );
}