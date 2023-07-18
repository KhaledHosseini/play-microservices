'use client'
import {ParseUserArray, User} from '../../types'
import React, { useEffect, useReducer, useState } from "react";
import { toast } from "react-hot-toast";

const UsersListComponent: React.FC = ()=> {
  
    const [isBusy,setIsBusy] = useState<boolean>(true)
    const [userList,setUserList] = useState<User[]>([])
    useEffect(() => {
        const fetchData = async () => {
          //TODO: add pagination
          const userResponse = await fetch("/api/user/list?page=1&size=50",{
            method: 'GET'
          });
          setIsBusy(false)
          const users: User[] = await ParseUserArray(userResponse)
          console.log("UsersList.useEffect.fetchData: Users array is: ",users)
          toast.success("Numbr of returned users: "+ users.length.toString())
          setUserList(users)
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
              <td className="p-2"> <div className="text-left font-medium text-gray-800">{user.name}</div> </td>
              <td className="p-2"> <div className="text-left font-medium text-green-500">{user.email}</div> </td>
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