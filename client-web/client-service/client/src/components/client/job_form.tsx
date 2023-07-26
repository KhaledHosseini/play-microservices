import React, { useEffect } from 'react'
import { useForm } from "react-hook-form";
import {Job, EmailJob, JobStatus, JobType} from '../../types'
interface JobFormProps {
  job: Job | null;
  onSubmit: (job: Job) => void;
  onClose:() => void;
}
const JobForm: React.FC<JobFormProps> = ({ job, onSubmit, onClose }) => {
  //
  const jobData: any | null = job?.JobData ? JSON.parse(job.JobData) : null
    const {
        handleSubmit,
        register,
        formState: { errors },
      } = useForm<EmailJob>({
        defaultValues: {
          Name: job?.Name,
          Description: job?.Description,
          ScheduleTime:job?.ScheduleTime,
          JobType: job?.JobType,
          JobData: job?.JobData,
          SourceAddress: jobData?.SourceAddress,
          DestinationAddress: jobData?.DestinationAddress,
          Subject: jobData?.Subject,
          Message: jobData?.Message
        },
      });
  
    const handleSubmitJob = (emailjob: EmailJob) => {
      const job: Job = {
        Id: '',
        Name: emailjob.Name,
        Description: emailjob.Description,
        ScheduleTime: emailjob.ScheduleTime,
        JobStatus: JobStatus.JOB_STATUS_UNKNOWN,
        JobType: JobType.JOB_TYPE_EMAIL,
        JobData: JSON.stringify({
          SourceAddress: emailjob.SourceAddress,
          DestinationAddress: emailjob.DestinationAddress,
          Subject: emailjob.Subject,
          Message: emailjob.Message
        })
      }
      onSubmit(job)
    };
  
    return (
        <form className="max-w-xl mx-auto bg-white p-8 rounded-md shadow-md"  onSubmit={handleSubmit(handleSubmitJob)}>
          <div className="grid gap-6 mb-6 lg:grid-cols-2">
            <div className="mb-4">
              <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="name">Name</label>
              <input className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                type="text" 
                id="name"
                placeholder="name"
                {...register('Name', { required: true })}/>
            </div>

            <div className="mb-4">
              <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="description">Description</label>
              <input className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                type="text" 
                id="description"
                placeholder="description"
                {...register('Description', { required: true })}/>
            </div>

            <div className="mb-4">
              <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="schedule_time">Schedule time</label>
              <input className='bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500'
                  id="dateTime"
                  type='text'
                  {...register('ScheduleTime') } />
            </div>
            <div className="mb-4">
              <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="job_type">Job type: 0 email, 1 sms</label>
              <input className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                type="number" 
                id="job_type"
                placeholder="0"
                {...register('JobType', { required: true })}/>
            </div>
          </div>
          
          <div className="flex items-start mb-6">
            <label htmlFor="job_data" className="ml-2 text-sm font-medium text-gray-900 dark:text-gray-400">Job Data</label>
          </div>

          <div className="grid gap-6 mb-6 lg:grid-cols-2">
            <div className="mb-4">
              <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="job_data.source_address">Source address</label>
              <input className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                type="email" 
                id="job_data.source_address"
                placeholder="admin@admin.com"
                {...register('SourceAddress', { required: true })}/>
            </div>
            <div className="mb-4">
              <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="job_data.destination_address">Destination address</label>
              <input className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                type="email" 
                id="job_data.destination_address"
                placeholder="customer@example.com"
                {...register('DestinationAddress', { required: true })}/>
            </div>
            <div className="mb-4">
              <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="job_data.subject">Subject</label>
              <input className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                type="text" 
                id="job_data.subject"
                placeholder="subject"
                {...register('Subject', { required: true })}/>
            </div>
            <div className="mb-4">
              <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="job_data.message">Message</label>
              <input className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                type="text" 
                id="job_data.message"
                placeholder="message"
                {...register('Message', { required: true })}/>
            </div>
          </div>

          <button className="w-full bg-indigo-500 text-white text-sm font-bold py-2 px-4 rounded-md hover:bg-indigo-600 transition duration-300" type="submit">Submit</button>
          <button onClick={()=>{onClose()}}  className="w-full bg-red-500 text-white text-sm font-bold py-2 px-4 rounded-md hover:bg-red-600 transition duration-300">Close</button>
        </form>
    );
  };
  
  export default JobForm;