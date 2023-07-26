import { StatusCodes } from 'http-status-codes';
// we do fetch, If the result is unauthorized, Possibily our access token has exired! we simply refresh it.
export async function api_fetch_with_access_token(url: string, options?: RequestInit | undefined): Promise<Response> {
    let result = await fetch(url,options)
    console.log("fetch result status code and texts are: ",result.status, result.statusText)
    if (result.status == StatusCodes.UNAUTHORIZED) {
        //possibly expired access token.
        console.log("Unauthorized. possibly access token has expired. lets refresh access token....")
        result = await fetch('/api/user/refresh_token', {
            method: "POST"
        })
        console.log("refresh token result arrived",result)
        if (result.status == StatusCodes.OK) {
            console.log("refresh token request is successfull. calling again.")
            //we have accessed new access token in browser cookie. now request again with our new access token
            result = await fetch(url,options)
            return result
        }else {
            console.log("refresh token request was not successfull")
        }
    }
    return result
}