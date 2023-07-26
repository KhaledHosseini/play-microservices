'use client'

import React, { useEffect, useState } from "react"
import { ParseUser } from "@/types";
import { useAuth } from "@/components/providers/auth_provider";
import Header from '@/components/client/header'

import UsersListComponent from '@/components/client/users_list_component'
import ReportListComponent from '@/components/client/report_list_component'
import JobListComponent from '@/components/client/job_list_component'
import {fetch_with_refresh_token} from '../../lib/api_gateway'
export default function Home() {

  const [_user,setUser] = useAuth()

  useEffect(() => {
    // declare the data fetching function
    const fetchData = async () => {
      const userResponse = await fetch_with_refresh_token("/api/user/get");
      const user = await ParseUser(userResponse)
      console.log("Home.useEffect.fetchData: User is: ",user)
      setUser(user)
    }
  
    fetchData()
    .catch(console.error);
  }, [setUser])

  const [currentComponent, setCurrentComponent] = useState<string | null>(null);

  return (
    <div>
      <Header/>
      <div className="flex min-h-screen flex-row bg-gray-100 text-gray-800">
        <aside className="sidebar w-48 -translate-x-full transform bg-white p-4 transition-transform duration-150 ease-in md:translate-x-0 md:shadow-md">
          <div className="my-4 w-full border-b-4 border-indigo-100 text-center">
            <div className="bg-white  rounded-md">
              <div className="bg-white rounded-md list-none  text-center ">
                <li className="py-3 border-b-2"><button onClick={()=>{setCurrentComponent('users')}} className="list-none  hover:text-indigo-600">Users</button></li>
                <li className="py-3 border-b-2"><button onClick={()=>{setCurrentComponent('jobs')}} className="list-none  hover:text-indigo-600">Scheduled Jobs</button></li>
                <li className="py-3 border-b-2"><button onClick={()=>{setCurrentComponent('reports')}} className="list-none  hover:text-indigo-600">Reports</button></li>
              </div>
            </div>
          </div>
        </aside>
        <main className="main -ml-48 flex flex-grow flex-col p-4 transition-all duration-150 ease-in md:ml-0">
          <div className="flex h-full items-start justify-center bg-white text-center text-5xl font-bold shadow-md">
            {currentComponent === 'users' && <UsersListComponent/>}
            {currentComponent === 'jobs' && <JobListComponent/>}
            {currentComponent === 'reports' && <ReportListComponent/>}
          </div>
        </main>
      </div>
    </div>
  );
}
