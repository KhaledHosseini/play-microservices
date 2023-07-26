'use client'
import {ParseJobArray, Job, CreateJob, JobType, JobStatus} from '../../types'
import React, { useEffect, useState } from "react";
import { toast } from "react-hot-toast";
import { useMutation } from '@tanstack/react-query';
import JobForm from './job_form';
import {api_fetch_with_access_token} from '../../lib/api_gateway'

async function deleteJob(id: string): Promise<Response> {
  const response = await api_fetch_with_access_token(`/api/job/delete?id=${id}`, {
    method: 'POST',
    headers:{},
    body:null
  })
  return response
}

async function createOrUpdateJob(job: Job,action: string): Promise<Response> {
  console.log("createOrUpdateJob started for action and job:",action,job)
  const url = action == "create" ? "/api/job/create" : "/api/job/update"
  var json = JSON.stringify(job)
  if (action == "create") {
    let createJob: CreateJob = {
      Name: job.Name,
      Description: job.Description,
      ScheduleTime: job.ScheduleTime,
      JobType: job.JobType,
      JobData: job.JobData
    }
    json = JSON.stringify(createJob)
  }
  const response = await api_fetch_with_access_token(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: json,
  })
  console.log("createOrUpdateJob returned with response",response)
  return response
}

const JobListComponent: React.FC = ()=> {

    const emptyJob: Job = {
      Id: '',
      Name: 'New job',
      Description: '',
      ScheduleTime: '',
      JobStatus: JobStatus.JOB_STATUS_UNKNOWN,
      JobType: JobType.JOB_TYPE_EMAIL,
      JobData: JSON.stringify({
        SourceAddress: '',
        DestinationAddress: '',
        Subject: '',
        Message: ''
      })
    }

    const [isBusy, setIsBusy] = useState<boolean>(true)
    const [jobList, setJobList] = useState<Job[]>([])
    useEffect(() => {
        // declare the data fetching function
        const fetchData = async () => {
          const jobResponse = await api_fetch_with_access_token("/api/job/list?page=1&size=100",{
            method: 'GET'
          });
          setIsBusy(false)
          console.log("response is ", jobResponse)
          if(jobResponse.ok){
            const jobs: Job[] = await ParseJobArray(jobResponse)
            console.log("ListJobComponent.useEffect.fetchData: job array is: ",jobs)
            toast.success("Numbr of returned jobs: "+ jobs.length.toString())
            setJobList(jobs)
          }else{
            toast.error("error:" + jobResponse.statusText)
          }
        }
      
        setIsBusy(true)
        fetchData()
        .catch(error => {
            toast.error(error.toString())
        });
      }, [])

      const deleteMutation = useMutation((id: string) => deleteJob(id), {
        onSuccess: (response, id) => {
          setIsBusy(false)
          if (response.status == 200){
            toast.success("delete success")
            const filteredArray = jobList.filter((job: Job) => job.Id !== id);
            setJobList(filteredArray)
          }else {
            toast.error("error deleting job")
          }
        },
        onError: ()=> {
          setIsBusy(false)
          toast.error("error deleting job")
        }
      });
  const handleDelete = (id: string) => {
    setIsBusy(true)
    deleteMutation.mutate(id);
  };

  const [jobForm_job, setJobForm_job] = useState<Job | null>(null)
  const [jobForm_action, setJobForm_action] = useState<string>('update')
  const createOrUpdateMutation = useMutation((job: Job) => createOrUpdateJob(job, jobForm_action), {
    onSuccess: (response, job) => {
      setIsBusy(false)
      if (response.status == 200){
        if (jobForm_action == "create") {
          jobList.push(job)
          setJobList(jobList)
        }else if (jobForm_action == "update") {
          const index = jobList.findIndex((obj: Job) => obj.Id === job.Id);
          if (index !== -1) {
              jobList.splice(index, 1, job);
              setJobList(jobList)
          }
        }
      }else {
        toast.error("error create or updating job")
      }
    },
    onError: ()=> {
      setIsBusy(false)
      toast.error("error create or updating job")
    }
  });

  const onSubmitForm = (job: Job) => {
    console.log("form submitted:")
    console.log(job)
    if (jobForm_job != null){
      job.Id = jobForm_job.Id; // if action is update.
    }
    createOrUpdateMutation.mutate(job)
    setJobForm_job(null)
  };

  const onCloseForm = () => {
    setJobForm_job(null)
  };

  return (
    <div className="container mx-auto px-4 sm:px-8">
      <div className="inline-flex items-center h-full ml-5 lg:w-2/5 lg:justify-end lg:ml-0">
        <button className='px-4 py-2 text-xs font-bold text-white uppercase transition-all duration-150 bg-teal-500 rounded shadow outline-none active:bg-teal-600 hover:shadow-md focus:outline-none ease'
        onClick={()=>{setJobForm_action("create");setJobForm_job(emptyJob);}}> Schedule New Job</button>
      </div>
    <div className="py-8">
    <div className="-mx-4 sm:-mx-8 px-4 sm:px-8 py-4 overflow-x-auto">
    <div className="inline-block min-w-full shadow rounded-lg overflow-hidden">
      {isBusy && <div className="text-sm text-gray-800">Busy...</div>}
      <table className="min-w-full">
        <thead className="bg-gray-50 text-xs font-semibold uppercase text-gray-400">
          <tr>
            <th></th>
            <th className="p-2"> <div className="text-left font-semibold">Name</div> </th>
            <th className="p-2"> <div className="text-left font-semibold">Description</div> </th>
            <th className="p-2"> <div className="text-left font-semibold">Schedule Time</div> </th>
            <th className="p-2"> <div className="text-left font-semibold">Status</div> </th>
            <th className="p-2"> <div className="text-left font-semibold">Type</div> </th>
            <th className="p-2"> <div className="text-left font-semibold">Edit</div> </th>
            <th className="p-2"> <div className="text-left font-semibold">Delete</div> </th>
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-100 text-sm">
          {jobList.map((job, index) => (
            <tr key={index}>
              <td className="p-2"> <div className="text-left font-medium text-gray-800">{index + 1}</div> </td>
              <td className="p-2"> <div className="text-left font-medium text-gray-500">{job.Name}</div> </td>
              <td className="p-2"> <div className="text-left font-medium text-gray-800">{job.Description.substring(0,25) + "..."}</div> </td>
              <td className="p-2"> <div className="text-left font-medium text-green-500">{job.ScheduleTime}</div> </td>
              <td className="p-2"> <div className="text-left font-medium text-blue-500">{JobStatus[job.JobStatus]}</div> </td>
              <td className="p-2"> <div className="text-left font-medium text-pink-500">{JobType[job.JobType]}</div> </td>
              <td className="p-2"> <div className="flex justify-center"> <button className='text-blue-600 dark:text-blue-500 hover:underline' onClick={()=> {setJobForm_action('update'); setJobForm_job(job)}}> Edit</button></div></td>
              <td className="p-2"> <div className="flex justify-center">
                <button onClick={()=> handleDelete(job.Id)}>
                    <svg className="h-8 w-8 rounded-full p-1 hover:bg-gray-100 hover:text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                    </svg>
                </button>
              </div></td> 
            </tr>
          ))}
        </tbody>
      </table>
      {jobForm_job != null ? (
        <>
          <div className="justify-center items-center flex overflow-x-hidden overflow-y-auto fixed inset-0 z-50 outline-none focus:outline-none">
            <JobForm job={jobForm_job} onSubmit={onSubmitForm } onClose={onCloseForm}/>
          </div>
          <div className="opacity-25 fixed inset-0 z-40 bg-black"></div>
        </>
      ) : null}
    </div>
    </div>
    </div>
    </div>
  );
}

export default JobListComponent