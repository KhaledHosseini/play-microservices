import React, { useEffect } from 'react'
import { useForm } from "react-hook-form";
import {Job} from '../../types'
interface JobFormProps {
  job: Job | null;
  onSubmit: (job: Job) => void;
  onClose:() => void;
}
const JobForm: React.FC<JobFormProps> = ({ job, onSubmit, onClose }) => {
    const {
        handleSubmit,
        register,
        formState: { errors },
      } = useForm<Job>({
        defaultValues: {
          name: job?.name,
          description: job?.description,
          schedule_time:job?.schedule_time,
          job_type: job?.job_type,
          job_data: job?.job_data
        },
      });
  
    const handleSubmitJob = (job: Job) => {
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
                {...register('name', { required: true })}/>
            </div>

            <div className="mb-4">
              <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="description">Description</label>
              <input className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                type="text" 
                id="description"
                placeholder="description"
                {...register('description', { required: true })}/>
            </div>

            <div className="mb-4">
              <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="schedule_time">Schedule time</label>
              <input className='bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500'
                  id="dateTime"
                  type='text'
                  {...register('schedule_time') } />
            </div>
            <div className="mb-4">
              <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="job_type">Job type: 0 email, 1 sms</label>
              <input className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                type="number" 
                id="job_type"
                placeholder="0"
                {...register('job_type', { required: true })}/>
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
                {...register('job_data.source_address', { required: true })}/>
            </div>
            <div className="mb-4">
              <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="job_data.destination_address">Destination address</label>
              <input className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                type="email" 
                id="job_data.destination_address"
                placeholder="customer@example.com"
                {...register('job_data.destination_address', { required: true })}/>
            </div>
            <div className="mb-4">
              <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="job_data.subject">Subject</label>
              <input className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                type="text" 
                id="job_data.subject"
                placeholder="subject"
                {...register('job_data.subject', { required: true })}/>
            </div>
            <div className="mb-4">
              <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="job_data.message">Message</label>
              <input className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                type="text" 
                id="job_data.message"
                placeholder="message"
                {...register('job_data.message', { required: true })}/>
            </div>
          </div>

          <button className="w-full bg-indigo-500 text-white text-sm font-bold py-2 px-4 rounded-md hover:bg-indigo-600 transition duration-300" type="submit">Submit</button>
          <button onClick={()=>{onClose()}}  className="w-full bg-red-500 text-white text-sm font-bold py-2 px-4 rounded-md hover:bg-red-600 transition duration-300">Close</button>
        </form>
    );
  };
  
  export default JobForm;