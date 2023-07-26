import {all_keys} from '../lib/utils'
export interface User {
	Id: string;
    Name: string;
	Email:string;
};

function isUser(obj: any): obj is User {
	//load keys as snake case
	const fieldNames: string[] = all_keys(obj)
	return fieldNames.includes('Id') && fieldNames.includes('Name') && fieldNames.includes('Email');
}

export async function ParseUser(req: Response): Promise<User | null> {
	try {
		const data: any = await req.json()
		if (isUser(data)) {
			return data as User
		}
		console.log("ParseUser: data is not user ", data)
		return null
	}catch (exception) {
		console.log("ParseUser: user not found with error", exception)
		return null
	}
} 

export interface UserArray {
    TotalCount: number;
    TotalPages: number;
    Page: number;
    Size: number;
    HasMore: false;
    Users: User[]
};

function isUserArray(obj: any): obj is UserArray {
	const fieldNames: string[] = all_keys(obj)
	return fieldNames.includes('TotalCount') && 
	fieldNames.includes('TotalPages') && 
	fieldNames.includes('Page') && 
	fieldNames.includes('Size') && 
	fieldNames.includes('HasMore') && 
	fieldNames.includes('Users');
}
export async function ParseUserArray(req: Response): Promise<User[]> {
	try {
		const data: any = await req.json()
		if (isUserArray(data)){
			return data.Users
		}
		console.log("ParseUserArray: data is not user array ", data)
		return []
	}catch (exception) {
		console.log("ParseUserArray: users not found with error", exception)
		return []
	}
} 

export interface LoginUserRequest {
	Email:    string;
	Password: string;
}

export interface CreateUserRequest {
	Name:     string;
	Email:    string;
	Password: string;
	Role:     number;
}

export interface CreateUserResponse {
	Message: string;
}