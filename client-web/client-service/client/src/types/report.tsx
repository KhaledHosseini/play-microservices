import {all_keys} from '../lib/utils'
export interface Report {
	Id: string;
	Topic:string;
    CreatedTime: string;
    ReportData: string;
};


export interface ReportArray {
    TotalCount: number;
    TotalPages: number;
    Page: number;
    Size: number;
    HasMore: false;
    Reports: Report[]
};

function isReport(obj: any): obj is Report {
	const fieldNames: string[] = all_keys(obj)
	return fieldNames.includes('Id') &&
	fieldNames.includes('Topic') && 
	fieldNames.includes('CreatedTime') && 
	fieldNames.includes('ReportData');
}

function isReportArray(obj: any): obj is ReportArray {
	const fieldNames: string[] = all_keys(obj)
	return fieldNames.includes('TotalCount') && 
	fieldNames.includes('TotalPages') && 
	fieldNames.includes('Page') && 
	fieldNames.includes('Size') && 
	fieldNames.includes('HasMore') && 
	fieldNames.includes('Reports');
}

export async function ParseReport(req: Response): Promise<Report | null> {
	try {
		const data = await req.json()
		if(isReport(data)){
			return data as Report;
		}
		console.log("ParseReport: data is not report", data)
		return null
	}catch (exception) {
		console.log("ParseReport: report not found with error", exception)
		return null
	}
} 

export async function ParseReportArray(req: Response): Promise<Report[]> {
	try {
        console.log("report request:",req)
		const data: any = await req.json()
		if (isReportArray(data)){
			return data.Reports
		}
		console.log("ParseReportArray: data is not report array ", data)
		return []
	}catch (exception) {
		console.log("ParseReportArray: reports not found with error", exception)
		return []
	}
}