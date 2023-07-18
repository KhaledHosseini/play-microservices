export interface Job {
    id: string;
    name: string;
    description: string;
    schedule_time: string;
    job_status: JobStatus;
    job_type: JobType;
    job_data: JobData;
};

export interface JobData {
    source_address: string;
    destination_address: string;
    subject: string;
    message: string;
}

export interface JobArray {
    total_count: number;
    total_pages: number;
    page: number;
    size: number;
    has_more: false;
    jobs: Job[]
};

export enum JobStatus {
    JOB_STATUS_UNKNOWN = 0,
    JOB_STATUS_PENDING = 1,
    JOB_STATUS_SCHEDULED = 2,
    JOB_STATUS_RUNNING = 3,
    JOB_STATUS_SUCCEEDED = 4,
    JOB_STATUS_FAILED = 5,
};
  
export enum JobType {
    JOB_TYPE_EMAIL = 0,
    JOB_TYPE_SMS = 1
};

export async function ParseJob(req: Response): Promise<Job | null> {
	try {
		const data: Job = await req.json()
		return data
	}catch (exception) {
		console.log("ParseUser: job not found with error", exception)
		return null
	}
} 

export async function ParseJobArray(req: Response): Promise<Job[]> {
	try {
        console.log("job request:",req)
		const data = await req.json()
        console.log("job json body:",data)
        const jobArray: JobArray = data as JobArray
        console.log("jobarray is:",jobArray)
		return jobArray.jobs;
	}catch (exception) {
		console.log("ParseJobArray: reports not found with error", exception)
		return []
	}
} 