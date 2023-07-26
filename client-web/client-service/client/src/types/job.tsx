import {all_keys} from '../lib/utils'
export interface Job {
    Id: string;
    Name: string;
    Description: string;
    ScheduleTime: string;
    JobStatus: JobStatus;
    JobType: JobType;
    JobData: string;
};

export interface CreateJob {
    Name: string;
    Description: string;
    ScheduleTime: string;
    JobType: JobType;
    JobData: string;
};

export interface JobArray {
    TotalCount: number;
    TotalPages: number;
    Page: number;
    Size: number;
    HasMore: false;
    Jobs: Job[]
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

export interface EmailJob {
    Id: string;
    Name: string;
    Description: string;
    ScheduleTime: string;
    JobStatus: JobStatus;
    JobType: JobType;
    JobData: string;
    SourceAddress: string;
    DestinationAddress: string;
    Subject: string;
    Message: string;
};

function isJob(obj: any): obj is Job {
    const fieldNames: string[] = all_keys(obj)
    return fieldNames.includes('Id') && 
    fieldNames.includes('Name') && 
    fieldNames.includes('Description') && 
    fieldNames.includes('ScheduleTime')&& 
    fieldNames.includes('JobStatus') && 
    fieldNames.includes('JobType') &&
    fieldNames.includes('JobData');
}

function isJobArray(obj: any): obj is JobArray {
    const fieldNames: string[] = all_keys(obj)
	return fieldNames.includes('TotalCount') && 
    fieldNames.includes('TotalPages') && 
    fieldNames.includes('Page') && 
    fieldNames.includes('Size') && 
    fieldNames.includes('HasMore') && 
    fieldNames.includes('Jobs');
}

export async function ParseJob(req: Response): Promise<Job | null> {
	try {
		const data: any = await req.json()
        if(isJob(data)) {
            return data as Job
        }
        console.log("ParseJob: data is not job", data)
		return null
	}catch (exception) {
		console.log("ParseJob: job not found with error", exception)
		return null
	}
} 

export async function ParseJobArray(req: Response): Promise<Job[]> {
	try {
        console.log("job request:",req)
		const data: any = await req.json()
        if(isJobArray(data)) {
            return data.Jobs
        }
        console.log("ParseJobArray: data is not job array", data)
		return []
	}catch (exception) {
		console.log("ParseJobArray: reports not found with error", exception)
		return []
	}
} 