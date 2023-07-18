import { Job } from "./job";
export interface Report {
	id: string;
    type: ReportType;
	topic:string;
    created_time: string;
    report_data: Job;
};

export enum ReportType {
	REPORT_TYPE_JOB = 0,
	REPORT_TYPE_UNKNOWN = 1
}

export interface ReportArray {
    total_count: number;
    total_pages: number;
    page: number;
    size: number;
    has_more: false;
    reports: Report[]
};

export async function ParseReport(req: Response): Promise<Report | null> {
	try {
		const data: Report = await req.json()
		return data
	}catch (exception) {
		console.log("ParseReport: report not found with error", exception)
		return null
	}
} 

export async function ParseReportArray(req: Response): Promise<Report[]> {
	try {
        console.log("report request:",req)
		const data = await req.json()
        console.log("report json body:",data)
        const reportArray: ReportArray = data as ReportArray
        console.log("report array is:",reportArray)
		return reportArray.reports;
	}catch (exception) {
		console.log("ParseReportArray: reports not found with error", exception)
		return []
	}
} 