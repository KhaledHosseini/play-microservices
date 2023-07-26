'use client'
import {ParseUserArray, User} from '../../types'
import React, { useEffect, useReducer, useState } from "react";
import { toast } from "react-hot-toast";
import {fetch_with_refresh_token} from '../../lib/api_gateway'
const UsersListComponent: React.FC = ()=> {
  
    const [isBusy,setIsBusy] = useState<boolean>(true)
    const [userList,setUserList] = useState<User[]>([])
    useEffect(() => {
        const fetchData = async () => {
          //TODO: add pagination
          const userResponse = await fetch_with_refresh_token("/api/user/list?page=1&size=100",{
            method: 'GET'
          });
          setIsBusy(false)
          console.log("response is ", userResponse)
          if(userResponse.ok){
            const users: User[] = await ParseUserArray(userResponse)
            console.log("UsersList.useEffect.fetchData: Users array is: ",users)
            toast.success("Numbr of returned users: "+ users.length.toString())
            setUserList(users)
          }else {
            toast.error("Error: "+ userResponse.statusText)
          }
        }
        setIsBusy(true)
        fetchData()
        .catch(error => {
            toast.error(error.toString())
        });
      }, [])

  return (
    <div className="container mx-auto px-4 sm:px-8">
    <div className="py-8">
    <div className="-mx-4 sm:-mx-8 px-4 sm:px-8 py-4 overflow-x-auto">
    <div className="inline-block min-w-full shadow rounded-lg overflow-hidden">
      {isBusy && <div className="text-sm text-gray-800">Busy...</div>}
      <table className="min-w-full">
        <thead className="bg-gray-50 text-xs font-semibold uppercase text-gray-400">
          <tr>
            <th></th>
            <th className="p-2"> <div className="text-left font-semibold">Name</div> </th>
            <th className="p-2"> <div className="text-left font-semibold">Email</div> </th>
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-100 text-sm">
          {userList.map((user, index) => (
            <tr key={index}>
              <td className="p-2"> <div className="text-left font-medium text-gray-800">{index + 1}</div> </td>
              <td className="p-2"> <div className="text-left font-medium text-gray-800">{user.Name}</div> </td>
              <td className="p-2"> <div className="text-left font-medium text-green-500">{user.Email}</div> </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
    </div>
    </div>
    </div>
  );
}

export default UsersListComponent