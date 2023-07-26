import { NextRequest, NextResponse } from 'next/server'
import {URL_APIGATEWAY} from './config'

export async function middleware(req: NextRequest) {
  const regex = new RegExp('/api/*')
  if (!regex.test(req.url)) {
    return new NextResponse(null, {
      status: 400,
      statusText: "Bad Request"
    })
  }
  console.log("middleware is called for url: ",req.url)
  const url = URL_APIGATEWAY + "/" + req.url.split("/api/")[1]
  console.log("middleware sends the request to : ", url)
  try {
    const res = await fetch( url, {
      method:req.method,
      headers: req.headers,
      body: req.body
    });
    return res
  }catch (err) {
    console.log(err)
    return new NextResponse('Internal Server Error', {
      status: 500,
      statusText: 'Internal Server(proxy) Error'
    });
  }
}

export const config = {
  matcher: '/api/:path*',
}