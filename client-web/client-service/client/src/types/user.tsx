export interface User {
	id: string;
    name: string;
	email:string;
};

export async function ParseUser(req: Response): Promise<User | null> {
	try {
		const data: User = await req.json()
		return data
	}catch (exception) {
		console.log("ParseUser: user not found with error", exception)
		return null
	}
} 

export async function ParseUserArray(req: Response): Promise<User[]> {
	try {
		const data = await req.json()
		return data.map((item: any) => item as User);
	}catch (exception) {
		console.log("ParseUser: user not found with error", exception)
		return []
	}
} 

export interface LoginUserRequest {
	email:    string;
	password: string;
}

export interface CreateUserRequest {
	name:     string;
	email:    string;
	password: string;
	role:     number;
}

export interface CreateUserResponse {
	message: string;
}