'use client'
import {ParseReportArray, Report, ReportType} from '../../types'
import React, { useEffect, useState } from "react";
import { toast } from "react-hot-toast";

const ReportListComponent: React.FC = ()=> {
    const [isBusy,setIsBusy] = useState<boolean>(true)
    const [errorMessage,setErrorMessage] = useState<string | null>(null)
    const [reportList,setReportList] = useState<Report[]>([])
    useEffect(() => {
        // declare the data fetching function
        const fetchData = async () => {
          const reportResponse = await fetch("/api/report/list?filter=1&page=1&size=10",{
            method: 'GET'
          });
          setIsBusy(false)
          const reports: Report[] = await ParseReportArray(reportResponse)
          console.log("UsersReportComponent.useEffect.fetchData: reports array is: ",reports)
          toast.success("Numbr of returned reports: "+ reports.length.toString())
          setReportList(reports)
        }
        setIsBusy(true)
        fetchData()
        .catch(error => {
            setErrorMessage(error)
        });
      }, [])

  return (
    <div className="container mx-auto px-4 sm:px-8">
    <div className="py-8">
    <div className="-mx-4 sm:-mx-8 px-4 sm:px-8 py-4 overflow-x-auto">
    <div className="inline-block min-w-full shadow rounded-lg overflow-hidden">
      {isBusy && <div className="text-sm text-gray-800">Busy...</div>}
      {errorMessage != null && <div className="font-semibold text-red-800">errorMessage</div>}
      <table className="min-w-full">
        <thead className="bg-gray-50 text-xs font-semibold uppercase text-gray-400">
          <tr>
            <th></th>
            <th className="p-2"> <div className="text-left font-semibold">Id</div> </th>
            <th className="p-2"> <div className="text-left font-semibold">Type</div> </th>
            <th className="p-2"> <div className="text-left font-semibold">Topic</div> </th>
            <th className="p-2"> <div className="text-left font-semibold">Data</div> </th>
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-100 text-sm">
          {reportList.map((report, index) => (
            <tr key={index}>
              <td className="p-2"> <div className="text-left font-medium text-gray-800">{index + 1}</div> </td>
              <td className="p-2"> <div className="text-left font-medium text-gray-800">{report.id}</div> </td>
              <td className="p-2"> <div className="text-left font-medium text-green-500">{ReportType[report.type]}</div> </td>
              <td className="p-2"> <div className="text-left font-medium text-gray-800">{report.topic}</div> </td>
              <td className="p-2"> <div className="text-left font-medium text-green-500">{JSON.stringify(report.report_data).substring(0, 60) + "..."}</div> </td>
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

export default ReportListComponent